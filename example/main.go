package main

import (
	init "github.com/johnksv/docker-compose-wrapper"
)

func main() {
	init.Flags()

	init.DockerComposeWrapper("", nil)
}
