package wrapper

import (
	"flag"
	"log"

	"github.com/johnksv/docker-compose-wrapper/internal"
	"github.com/johnksv/docker-compose-wrapper/types"
)

var ComposeFile string
var Basedir string
var EnvFile string
var Registry string

func Flags() {
	composeFilePtr := flag.String("f", "", "The docker-compose file to use.")
	basedirPtr := flag.String("d", internal.GetWorkingDir(), "The directory to scan for compose and .env files. Defaults to working dir.")
	envFilePtr := flag.String("env", "", "The .env file to use.")
	registryPtr := flag.String("r", "", "The registry to use.")

	flag.Parse()
	ComposeFile = *composeFilePtr
	Basedir = *basedirPtr
	EnvFile = *envFilePtr
	Registry = *registryPtr
}

// DockerComposeWrapper Registry: The registry to use for pulling images. . EnvLookup: custom env lookup, can be nil
func DockerComposeWrapper(envLookup types.EnvLookup) {

	if ComposeFile == "" && EnvFile == "" {
		internal.Run(Basedir, Registry, envLookup)
	} else if ComposeFile != "" && EnvFile != "" {
		internal.RunWithFile(ComposeFile, EnvFile, Registry, envLookup)
	} else {
		log.Fatal("env-flag and f-flag must be used together")
	}
}
