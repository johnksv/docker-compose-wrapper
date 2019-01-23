package main

import (
	wrapper "github.com/johnksv/docker-compose-wrapper"
)

func main() {
	wrapper.Flags()

	wrapper.DockerComposeWrapper("", nil)
}
