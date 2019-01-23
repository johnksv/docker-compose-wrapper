package main

import (
	"docker-compose-wrapper/internal"
)

func main() {
	internal.InitFlags()

	internal.DockerComposeWrapper("", nil)
}
