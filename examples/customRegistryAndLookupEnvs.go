package main

import (
	"github.com/johnksv/docker-compose-wrapper/wrapper"
)

func main() {
	wrapper.InitFlags("custom.registry/")

	lookupEnv := func(key string, defaultValue string) string {
		if key == "LOOKUP_THIS_ENV" {
			return "value"
		}
		return defaultValue
	}

	wrapper.RunDockerComposeWrapper(lookupEnv)
}
