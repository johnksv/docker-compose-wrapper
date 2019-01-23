package internal

import (
	"docker-compose-wrapper/types"
	"fmt"
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
)

func findEnvFile(basedir string) string {
	matches, err := filepath.Glob(basedir + "/.env")
	if err != nil {
		log.Fatal(err)
	}
	if len(matches) == 0 {
		log.Fatal("No .env file found. Not supported yet.")
	}
	fmt.Println("Using .env-file at:", matches[0])
	return matches[0]
}

func parseEnvFile(envFile string) map[string]envVal {
	envVariables, err := godotenv.Read(envFile)
	if err != nil {
		log.Fatal(err)
	}

	result := make(map[string]envVal)
	for key, val := range envVariables {
		result[key] = envVal{
			value: val,
		}
	}

	return result
}

func usedEnvs(servicesToRun []Service, envsToCheck map[string]envVal) (map[string]envVal, []Service) {
	var result = map[string]envVal{}
	var resultServices []Service

	for _, service := range servicesToRun {
		usedEnvs := service.getUsedSystemEnvs()
		for _, key := range usedEnvs {
			_, ok := envsToCheck[key]
			if ok {
				result[key] = envsToCheck[key]
			} else if key != "" { // If key doesn't exists in .env (but it might need to be looked up)
				result[key] = envVal{value: ""}
			}
		}
		resultServices = append(resultServices, service)
	}
	return result, resultServices
}

func lookupEnvs(envLookup types.EnvLookup, envs map[string]envVal) {
	for key, envValue := range envs {
		envs[key] = envVal{
			value: envLookup(key, envValue.value),
		}
	}
}
