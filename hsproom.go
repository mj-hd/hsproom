package main

import (
	"flag"

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
			log.Fatal(err)
		}

	}

	gum.Start()

}
