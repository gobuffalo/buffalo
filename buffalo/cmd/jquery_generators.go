package cmd

import "github.com/markbates/gentronics"

func newJQueryGenerator() *gentronics.Generator {
	should := func(data gentronics.Data) bool {
		if p, ok := data["withJQuery"]; ok {
			return p.(bool)
		}
		return false
	}

	g := gentronics.New()
	jf := &gentronics.RemoteFile{
		File: gentronics.NewFile("assets/jquery.js", ""),
	}
	jf.Should = should
	jf.RemotePath = "https://cdnjs.cloudflare.com/ajax/libs/jquery/3.1.1/jquery.min.js"
	g.Add(jf)

	jm := &gentronics.RemoteFile{
		File: gentronics.NewFile("assets/jquery.map", ""),
	}
	jm.Should = should
	jm.RemotePath = "https://cdnjs.cloudflare.com/ajax/libs/jquery/3.1.1/jquery.min.map"
	g.Add(jm)
	return g
}
