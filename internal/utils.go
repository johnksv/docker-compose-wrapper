package internal

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func arrayToString(arr []string) string {
	return strings.Join(arr, " ")
}

func getWorkingDir() string {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return workingDir
}

func printConfiguration(services []Service, usedEnvs map[string]envVal, registry string) {
	fmt.Println("----------------------------------------------------------------")
	format := "%-20s%-50s%-20s"
	fmt.Printf("%-20s%-50s%-20s%-10s", "Registry", "Service", "Version", "Other")
	fmt.Println()

	for _, service := range services {
		imageEnv := service.getImageEnv(usedEnvs, registry)

		version := imageEnv.version
		if version == "" {
			version = "Will build on run"
		}

		fmt.Printf(format, imageEnv.registry, imageEnv.image, version)

		for _, key := range service.getUsedSystemEnvs() {
			if !strings.Contains(strings.ToLower(key), "image") {
				valueString := usedEnvs[key].value
				if len(valueString) == 0 {
					valueString = "\"\""
				}
				fmt.Print(key, "=", valueString)
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func envPropertiesToString(arr map[string]string) string {
	var builder strings.Builder
	for key, value := range arr {
		builder.WriteString(key + "=" + value)
		builder.WriteRune(' ')
	}

	return strings.Trim(builder.String(), " ")
}
