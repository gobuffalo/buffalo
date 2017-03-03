package assets

import "github.com/markbates/gentronics"

const LogoURL = "https://raw.githubusercontent.com/gobuffalo/buffalo/master/logo.svg"

var PublicLogo = &gentronics.RemoteFile{
	File:       gentronics.NewFile("public/assets/images/logo.svg", ""),
	RemotePath: LogoURL,
}

var AssetsLogo = &gentronics.RemoteFile{
	File:       gentronics.NewFile("assets/images/logo.svg", ""),
	RemotePath: LogoURL,
}
