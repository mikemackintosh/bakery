Bakery
-----------
[![Actions Status](https://github.com/mikemackintosh/bakery/workflows/Test/badge.svg)](https://github.com/mikemackintosh/bakery/actions)

<p align="center">
  <img width="128px" src=".github/bakery.png">
</p>

Bakery is a sweet new way of configuring devices with configuration formats you're familiar with.

### Overview
I lead a client engineering team for a few years. We were very effective at our jobs, but often were annoyed at the inconsistent tooling between macOS, Windows and Linux. macOS required munki or worse off, JAMF. Linux ran Chef, which became costly in an enterprise due to their changing business model. Windows, no one wanted to tackle that giant leaving unstable and unpolished tools like AirWatch part of our limited toolset. We knew there had to be a better way.

Adding onto a tool I used personally for several years to manage my own devices, I wanted to make a proof of concept for something better just to show it's possible. I sure as hell appreciate and respect all the work out there to date; all the tools and services that other contributors have shared and maintained.

I started with macOS support due to it being my daily driver.

#### Goals and Values
Keeping Bakery simplistic is the number one goal, and it's easy to stray from this. In order to stay focused, I came up with the following values:

  - The framework you use to deliver your configurations should have a small physical footprint on the client.
  - The framework delivering the configuration should have a minimal number of runtime dependencies.
  - The framework should be able to use a bundled config, a remote config or a local config.
  - The configurations should be easily readable.
  - The configurations should be easily writeable.
  - The configurations should be easily testable.

By defining the above, I came to the following conclusions:
  - I wanted the service to run in Go, since there are no runtime dependencies. Once the binary is built, it can run without additional software.
  - I wanted the service to be able to pack local configs as part of the binary.
  - I want to use HCL as the standard dsl since it's very clear, explicit and familiar to most engineers.
  - I wanted the service to be able to download remote configurations.
  - I wanted the service to be able to bundle package resources.

## Usage
Download Bakery:

    go get github.com/mikemackintosh/bakery
    cd $GOPATH/src/github.com/mikemackintosh/bakery

Build:

    make build

Update your configuration:

    vim config.yum

### Flags:

    Usage of bakery:
      --temp-dir string
        	Temporary resource directory (default " /var/bakery/tmp")
      -b	Bundle client config with binary
      -c string
        	Configuration file (default "manifest.yml")
      -d	When enabled, turns on debugging
      -r string
        	Client recipe file (default "config.yum")
      -v int
        	Sets output verbosity level (default 1)

## Resource Types
The following are just a preview of resource types supported. There is also dependency resolution which you will see in the examples below.

#### Git
```
git "dotfiles" {
  source = "https://github.com/mikemackintosh/dotfiles"
  destination = "~/.dotfiles_test"
  branch = "master"
  depends_on = "Accept Xcode License"
  path = "/usr/bin/git"
  user = "self"
}
```

#### Shell
```
shell "Accept Xcode License" {
  script = <<EOT
xcodebuild -license accept
EOTs
}
```

#### Brew
```
brew "package" {
  action = "upgrade"
}
```

#### Zip
```
// Install dash from their website (.app bundle within the zip)
zip "Dash" {
  source = "https://sanfrancisco.kapeli.com/downloads/v4/Dash.zip"
  checksum = "802c5a63ac72c94ae4c6481529f11795c35f316d3607c9946f777f447b670c50"
  destination = "/Applications/"
  not_if = "ls -la /Applications/ | grep 'Dash.app'"
}
```


#### DMG
```
dmg "Docker" {
  source = "https://download.docker.com/mac/stable/Docker.dmg"
  checksum = "a06307d8da9c3778b183786cb87037ed3b7226d36ebc978fd40aa90851c0a04e"
}
```

#### Font
```
font "ubuntu" {
  source = "https://assets.ubuntu.com/v1/fad7939b-ubuntu-font-family-0.83.zip"
  checksum = "456d7d42797febd0d7d4cf1b782a2e03680bb4a5ee43cc9d06bda172bac05b42"
}
```
