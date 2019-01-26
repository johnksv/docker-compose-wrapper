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
var registry string

func InitFlags(defaultRegistry string) {
	composeFilePtr := flag.String("f", "", "The docker-compose file to use.")
	basedirPtr := flag.String("d", internal.GetWorkingDir(), "The directory to scan for compose and .env files. Defaults to working dir.")
	envFilePtr := flag.String("env", "", "The .env file to use.")
	registryPtr := flag.String("r", defaultRegistry, "The registry to use.")

	flag.Parse()
	ComposeFile = *composeFilePtr
	Basedir = *basedirPtr
	EnvFile = *envFilePtr
	registry = *registryPtr
}

// DockerComposeWrapper. EnvLookup: custom env lookup, can be nil
func RunDockerComposeWrapper(envLookup types.EnvLookup) {

	if ComposeFile == "" && EnvFile == "" {
		internal.Run(Basedir, registry, envLookup)
	} else if ComposeFile != "" && EnvFile != "" {
		internal.RunWithFile(ComposeFile, EnvFile, registry, envLookup)
	} else {
		log.Fatal("env-flag and f-flag must be used together")
	}
}
