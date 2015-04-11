package main

import (
	"flag"
	"os"

	"hsproom/gum"
	"hsproom/utils"
)

func main() {

	defer gum.Del()

	noDaemonize := flag.Bool("nodaemonize", false, "Do not daemonize")
	flag.Parse()

	if !*noDaemonize {

		err := gum.Daemonize()

		if err != nil {
			utils.PromulgateFatal(os.Stdout, err)
		}

	}

	gum.Start()

}
