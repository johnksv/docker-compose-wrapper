package internal

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	"github.com/johnksv/docker-compose-wrapper/types"
)

const (
	up   string = "up"
	pull string = "pull"
)

/*
Read compose
Read env
Select services
Lookup used envs
Config OK, mark which envs that are used
	Set used envs to custom value
Run commands
*/

type LookupEnv interface {
	lookupEnv(name string) string
}

func Run(basedir string, registry string, customLookupFunc types.EnvLookup) {
	composeFile := searchAndAskForDockerComposeFile(basedir)
	envFile := findEnvFile(basedir)
	// composeFile := "/Users/johnksv/digipost/digipost/docker-compose.yml"
	// envFile := "/Users/johnksv/digipost/digispost/.env"
	RunWithFile(composeFile, envFile, registry, customLookupFunc)
}

func RunWithFile(composeFile string, envFile string, registry string, customLookupFunc types.EnvLookup) {
	services := parse(composeFile)
	envFileValues := parseEnvFile(envFile)
	// defaultEnvValues := envPropertiesToString(envFileValues)

	// serviceAndEnvs := usedEnvVars(services, envFileValues)
	sortOnServiceName(services)

	servicesToRun := askWhichServicesToRun(services)
	if len(servicesToRun) == 0 {
		fmt.Println("No services chosen")
		os.Exit(1)
	}

	usedEnvs, servicesToRun := usedEnvs(servicesToRun, envFileValues)

	if customLookupFunc != nil {
		lookupEnvs(customLookupFunc, usedEnvs)
	}

	for {
		printConfiguration(servicesToRun, usedEnvs, registry)
		if askOneConfirm("Configuration OK?") {
			break
		} else {
			changeConfig(servicesToRun, usedEnvs, registry)
		}
	}

	var commands []string
	if askOneConfirm("Pull from registry (overwrite if old container with same tag exists)?") {
		commands = append(commands, pull)
	}

	commands = append(commands, up)
	for _, cmd := range commands {
		runDockerCompose(composeFile, cmd, usedEnvs, servicesToRun, registry)
	}

	fmt.Println("To stop, use: docker-compose stop")
}

func sortOnServiceName(services []Service) {
	sort.Slice(services, func(i, j int) bool { return services[i].Name < services[j].Name })
}

func runDockerCompose(dockerComposeFile string, action string, additionalEnvs map[string]envVal, services []Service, registry string) {

	var envs = os.Environ()
	for key, envVal := range additionalEnvs {
		envs = append(envs, key+"="+envVal.value)
	}

	args := []string{"-f", dockerComposeFile, action}

	var anyServices bool
	for _, service := range services {
		if action == pull && service.getImageEnv(additionalEnvs, registry).registry == "" {
			continue
		}
		args = append(args, service.Name)
		anyServices = true
	}
	if action == pull && !anyServices {
		return
	}

	cmd := exec.Command("docker-compose", args...)
	cmd.Env = envs
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Running:", cmd.Args)
	err := cmd.Run()
	if err != nil {
		log.Fatal("Error running docker-compose:", err)
	}
}

func searchAndAskForDockerComposeFile(basedir string) string {
	matches, err := filepath.Glob(basedir + "/docker-compose*.yml")
	if err != nil {
		log.Fatal(err)
	}
	if len(matches) == 0 {
		log.Fatal("No docker-compose found. Run with flag -h for help")
	}

	result := matches[0]

	if len(matches) > 1 {
		askWhichFileToUse(matches, &result)
	}

	return result
}
