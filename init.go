package main

import (
	"docker-compose-wrapper/internal"
	"docker-compose-wrapper/types"
	"flag"
	"log"
)

var composeFile string
var basedir string
var envFile string

func InitFlags() {
	composeFilePtr := flag.String("f", "", "The docker-compose file to use.")
	basedirPtr := flag.String("d", internal.GetWorkingDir(), "The directory to scan for compose and .env files. Defaults to working dir.")
	envFilePtr := flag.String("env", "", "The .env file to use.")

	flag.Parse()
	composeFile = *composeFilePtr
	basedir = *basedirPtr
	envFile = *envFilePtr
}

// DockerComposeWrapper Registry: The registry to use for pulling images. EnvLookup: custom env lookup, can be nil
func DockerComposeWrapper(registry string, envLookup types.EnvLookup) {

	if composeFile == "" && envFile == "" {
		internal.Run(basedir, registry, envLookup)
	} else if composeFile != "" && envFile != "" {
		internal.RunWithFile(composeFile, envFile, registry, envLookup)
	} else {
		log.Fatal("env-flag and f-flag must be used together")
	}
}
