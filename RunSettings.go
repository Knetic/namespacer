package main

import (
	"flag"
)

type RunSettings struct {
	targetPath string
	namespace string
}

func parseRunSettings() (RunSettings, error) {

	var ret RunSettings

	flag.Parse()
	ret.targetPath = flag.Arg(0)
	ret.namespace = flag.Arg(1)
	
	return ret, nil
}
