package main

import (
	"github.com/johnksv/docker-compose-wrapper/wrapper"
)

func main() {
	wrapper.InitFlags("")

	wrapper.RunDockerComposeWrapper(nil)
}
