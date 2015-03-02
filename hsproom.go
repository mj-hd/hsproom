package main

import (
	"os"
	"hsproom/gum"
	"hsproom/utils"
)

func main() {

	defer gum.Del()

	err := gum.Daemonize()

	if err != nil {
		utils.PromulgateFatal(os.Stdout, err)
	}

	gum.Start()

}
