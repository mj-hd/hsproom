package main

import (
	"./gum"
)

func main() {

	defer gum.Del()

	gum.Start()

}
