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
	AcceptEula     bool     `json:"accept_eula"`
	AllowUntrusted bool     `json:"allow_untrusted"`
	Checksum       *string  `json:"checksum"`
	App            *string  `json:"app"`
	Destination    *string  `json:"destination"`
	Source         string   `json:"source"`
	DependsOn      []string `json:"depends_on"`
}

// identifies the DMG spec
var dmgSpec = &hcldec.ObjectSpec{
	"depends_on": dependsOn,
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
}

// GetDestination will get or return the default destination
func (d *Dmg) GetDestination() string {
	if d.Destination != nil {
		destination := *d.Destination
		if string(destination[len(destination)-1]) != "/" {
			destination = destination + "/"
		}
		return destination
	}

	return "/Applications/"
}

// Parse will parse the config with the spec
func (d *Dmg) Parse(evalContext *hcl.EvalContext) error {
	cli.Debug(cli.INFO, "Preparing DMG", d.Name)
	cfg, diags := hcldec.Decode(d.Config, dmgSpec, evalContext)
	if len(diags) != 0 {
		for _, diag := range diags {
			cli.Debug(cli.INFO, "\t#", diag)
		}
		return fmt.Errorf("%s", diags.Errs()[0])
	}

	err := d.Populate(cfg, d)
	if err != nil {
		return err
	}

	return nil
}

// Bake will perform the DMG installation
func (d *Dmg) Bake() {
	// Rsync the app from the mounted DMG to the destination folder
	var appName = d.Name
	if d.App != nil {
		appName = *d.App
	}
	// Sets the app name with the .app extension
	var appNameWithExt = appName + ".app"

	u, err := url.Parse(d.Source)
	if err != nil {
		cli.Debug(cli.INFO, fmt.Sprintf("Error finding source %s", d.Source), err)
	}

	var tmpFile string
	if u.Scheme == "http" || u.Scheme == "https" {
		cli.Debug(cli.DEBUG, "\t-> Using HTTP(s) source for download", nil)

		urlParse, urlErr := url.Parse(d.Source)
		if err != nil {
			cli.Debug(cli.INFO, fmt.Sprintf("Error finding source %s", d.Source), urlErr)
		}

		tmpFile = config.Registry.TempDir + "/" + path.Base(urlParse.Path)
		urlErr = DownloadFile(d.Source, tmpFile, d.Checksum)
		if urlErr != nil {
			cli.Debug(cli.INFO, fmt.Sprintf("Error downloading file %s", d.Source), urlErr)
		}
	}

	// Mount it
	var mountpoint = fmt.Sprintf("/Volumes/%s", d.Name)
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
		d.GetDestination()}
	cli.Debug(cli.INFO, fmt.Sprintf("Installing %s to %s", tmpFile, d.GetDestination()), err)
	cli.Debug(cli.DEBUG2, fmt.Sprintf("\t-> Install command: %s", strings.Join(installCmd, " ")), err)
	r, err = RunCommand(installCmd)
	if err != nil {
		cli.Debug(cli.INFO, fmt.Sprintf("Error installing %s", appName), err)
	}
	cli.Debug(cli.DEBUG2, fmt.Sprintf("\t-> Install command response: \n%s", r.FormattedString()), err)

	// unmount the DMG after copying it over
	_, err = RunCommand([]string{hdiutilBinary, "unmount", mountpoint})
	if err != nil {
		cli.Debug(cli.INFO, "Error unmounting", err)
	}
}
