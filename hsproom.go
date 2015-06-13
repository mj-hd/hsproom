package main

import (
	"flag"
	"os"

	"./gum"
	"./utils/log"
)

func main() {

	defer gum.Del()

	noDaemonize := flag.Bool("nodaemonize", false, "Do not daemonize")
	flag.Parse()

	if !*noDaemonize {

		err := gum.Daemonize()

		if err != nil {
			log.Fatal(os.Stdout, err)
		}

	}

	gum.Start()

}
