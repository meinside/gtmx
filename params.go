// params.go

package main

// parameter definitions
type params struct {
	PrintVersion       bool `short:"V" long:"version" description:"Print version"`
	GenerateConfig     bool `short:"g" long:"gen-config" description:"Print a sample config file to stdout"`
	ListSessions       bool `short:"l" long:"list" description:"List sessions"`
	QuitCurrentSession bool `short:"q" long:"quit" description:"Quit current session"`
	Verbose            bool `short:"v" long:"verbose"`
}

func (p *params) multipleTaskRequested() bool {
	requested := 0

	if p.PrintVersion {
		requested += 1
	}
	if p.GenerateConfig {
		requested += 1
	}
	if p.ListSessions {
		requested += 1
	}
	if p.QuitCurrentSession {
		requested += 1
	}

	return requested > 1
}
