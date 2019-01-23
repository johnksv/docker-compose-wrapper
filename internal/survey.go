package internal

import (
	"fmt"
	"os"
	"sort"
	"strings"

	survey "gopkg.in/AlecAivazis/survey.v1"
)

func askWhichFileToUse(options []string, surveyResult *string) {
	surevySelect := &survey.Select{
		Message:  "Which docker-compose file to use?",
		Options:  options,
		PageSize: 15,
	}
	askOne(surevySelect, surveyResult)
}

func askWhichServicesToRun(services []Service) []Service {
	var options []string
	for _, e := range services {
		options = append(options, e.Name)
	}
	multislect := &survey.MultiSelect{
		Message:  "Which applications do you want to run?",
		Options:  options,
		PageSize: 15,
	}

	var surveyResult []string
	askOne(multislect, &surveyResult)

	var result []Service
	for _, selectedService := range surveyResult {
		index := sort.Search(len(services), func(i int) bool {
			return string(services[i].Name) >= selectedService
		})
		result = append(result, services[index])
	}

	return result
}

func askOneConfirm(question string) bool {
	confirm := &survey.Confirm{
		Message: question,
		Default: true,
	}
	var answer bool
	askOne(confirm, &answer)
	return answer
}

func askOne(question survey.Prompt, answer interface{}) {
	err := survey.AskOne(question, answer, nil)

	if err != nil {
		os.Exit(1)
	}
}

func changeConfig(services []Service, usedEnvs map[string]envVal, registryDefault string) {
	question := func(message string, defaultVal string) *survey.Input {
		return &survey.Input{
			Message: message,
			Default: defaultVal,
		}
	}

	for _, service := range services {
		fmt.Println("---------------------------")
		fmt.Println("Config for", service.Name)

		imgProps := service.getImageEnv(usedEnvs, registryDefault)

		registry := ""
		if registryDefault != "" {
			pullFromRegistry := askOneConfirm("Pull from registry (" + registryDefault + ")?")
			if pullFromRegistry {
				registry = registryDefault
			}
		}

		version := ""
		askOne(question("Version", imgProps.version), &version)

		imgKey := strings.TrimPrefix(service.Image, "$")
		usedEnvs[imgKey] = envVal{
			value: registry + imgProps.image + ":" + version,
		}

		for _, key := range service.getUsedSystemEnvs() {
			if !strings.Contains(strings.ToLower(key), "image") {
				answer := ""
				askOne(question("Value "+key, usedEnvs[key].value), &answer)
				usedEnvs[key] = envVal{
					value: answer,
				}
			}
		}
	}
}
