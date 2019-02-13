package pantry

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcldec"
	"github.com/mikemackintosh/bakery/cli"
	"github.com/mikemackintosh/bakery/config"
	"github.com/zclconf/go-cty/cty"
)

// Dmg is a MacOS DMG object
type Dmg struct {
	PantryItem
	Name           string   `hcl:"name,label"`
	Config         hcl.Body `hcl:",remain"`
	App            *string  `json:"app"`
	Source         string   `json:"source"`
	Destination    *string  `json:"destination"`
	Checksum       *string  `json:"checksum"`
	AcceptEula     bool     `json:"accept_eula"`
	AllowUntrusted bool     `json:"allow_untrusted"`
	Force          bool     `json:"force"`
}

// identifies the DMG spec
var dmgSpec = NewPantrySpec(&hcldec.ObjectSpec{
	"source": &hcldec.AttrSpec{
		Name:     "source",
		Required: true,
		Type:     cty.String,
	},
	"app": &hcldec.AttrSpec{
		Name:     "app",
		Required: false,
		Type:     cty.String,
	},
	"checksum": &hcldec.AttrSpec{
		Name:     "checksum",
		Required: true,
		Type:     cty.String,
	},
	"accept_eula": &hcldec.AttrSpec{
		Name:     "accept_eula",
		Required: false,
		Type:     cty.Bool,
	},
	"destination": &hcldec.AttrSpec{
		Name:     "destination",
		Required: false,
		Type:     cty.String,
	},
	"allow_untrusted": &hcldec.AttrSpec{
		Name:     "allow_untrusted",
		Required: false,
		Type:     cty.Bool,
	},
	"force": &hcldec.AttrSpec{
		Name:     "force",
		Required: false,
		Type:     cty.Bool,
	},
})

// GetDestination will get or return the default destination
func (p *Dmg) GetDestination() string {
	if p.Destination != nil {
		destination := *p.Destination
		if string(destination[len(destination)-1]) != "/" {
			destination = destination + "/"
		}
		return destination
	}

	return "/Applications/"
}

// Parse will parse the config with the spec
func (p *Dmg) Parse(evalContext *hcl.EvalContext) error {
	cli.Debug(cli.INFO, "Preparing DMG", p.Name)
	cfg, diags := hcldec.Decode(p.Config, dmgSpec, evalContext)
	if len(diags) != 0 {
		for _, diag := range diags {
			cli.Debug(cli.INFO, "\t#", diag)
		}
		return fmt.Errorf("%s", diags.Errs()[0])
	}

	err := p.Populate(cfg, p)
	if err != nil {
		return err
	}

	return nil
}

// Bake will perform the DMG installation
func (p *Dmg) Bake() {
	// Rsync the app from the mounted DMG to the destination folder
	var appName = p.Name
	if p.App != nil {
		appName = *p.App
	}
	// Sets the app name with the .app extension
	var appNameWithExt = appName + ".app"

	if FileExists(p.GetDestination()+appNameWithExt) && !p.Force {
		cli.Debug(cli.INFO, "\t-> Package already exists", p.GetDestination()+appNameWithExt)
		return
	}

	u, err := url.Parse(p.Source)
	if err != nil {
		cli.Debug(cli.INFO, fmt.Sprintf("Error finding source %s", p.Source), err)
	}

	var tmpFile string
	if u.Scheme == "http" || u.Scheme == "https" {
		cli.Debug(cli.DEBUG, "\t-> Using HTTP(s) source for download", nil)

		urlParse, urlErr := url.Parse(p.Source)
		if err != nil {
			cli.Debug(cli.INFO, fmt.Sprintf("Error finding source %s", p.Source), urlErr)
		}

		tmpFile = config.Registry.TempDir + "/" + path.Base(urlParse.Path)
		urlErr = DownloadFile(p.Source, tmpFile, p.Checksum)
		if urlErr != nil {
			cli.Debug(cli.INFO, fmt.Sprintf("Error downloading file %s", p.Source), urlErr)
		}
	}

	// Mount it
	var mountpoint = fmt.Sprintf("/Volumes/%s", p.Name)
	var hdiutilBinary = "/usr/bin/hdiutil"
	var mountCmd = []string{
		hdiutilBinary,
		"attach",
		strings.Replace(tmpFile, " ", "\\ ", -1),
		"-nobrowse",
		"-mountpoint",
		mountpoint}

	cli.Debug(cli.INFO, fmt.Sprintf("Mounting %s", tmpFile), err)
	cli.Debug(cli.DEBUG2, fmt.Sprintf("\t-> Install command: %#v", strings.Join(mountCmd, " ")), err)
	r, err := RunCommand(mountCmd)
	if err != nil {
		cli.Debug(cli.ERROR, fmt.Sprintf("Error mounting %s to %s", tmpFile, mountpoint), err)
	}
	cli.Debug(cli.DEBUG2, fmt.Sprintf("\t-> Mount command response: \n%s", r.FormattedString()), err)

	var installCmd = []string{
		"sudo",
		"/usr/bin/rsync",
		"--force",
		"--recursive",
		"--links",
		"--perms",
		"--executability",
		"--owner",
		"--group",
		"--times",
		fmt.Sprintf("%s/%s", mountpoint, appNameWithExt),
		p.GetDestination()}
	cli.Debug(cli.INFO, fmt.Sprintf("Installing %s to %s", tmpFile, p.GetDestination()), err)
	cli.Debug(cli.DEBUG2, fmt.Sprintf("\t-> Install command: %s", strings.Join(installCmd, " ")), err)
	r, err = RunCommand(installCmd)
	if err != nil {
		cli.Debug(cli.INFO, fmt.Sprintf("Error installing %s", appName), err)
	}
	cli.Debug(cli.DEBUG2, fmt.Sprintf("\t-> Install command response: \n%s", r.FormattedString()), err)

	// unmount the DMG after copying it over
	var unmountCommand = []string{
		hdiutilBinary,
		"unmount",
		mountpoint,
	}
	_, err = RunCommand(unmountCommand)
	if err != nil {
		cli.Debug(cli.INFO, "Error unmounting", err)
	}
}
