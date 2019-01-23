package internal

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/docker/cli/cli/compose/loader"
)

func mapBuild(build interface{}) map[string]interface{} {
	if build == nil {
		return nil
	}
	return build.(map[string]interface{})
}

func mapVolumes(volumes interface{}) []volum {
	if volumes == nil {
		return nil
	}
	var result []volum

	for _, content := range volumes.([]interface{}) {
		contentString := content.(string)

		parts := strings.Split(contentString, ":")

		var permission string
		if len(parts) > 2 {
			permission = parts[2]
		}

		result = append(result, volum{
			localPath:  parts[0],
			monutPath:  parts[1],
			permission: permission,
		})

	}

	return result
}

func mapEnviornment(envsInterface interface{}) map[string]string {
	if envsInterface == nil {
		return nil
	}

	var result = map[string]string{}

	//List of map[string]string, but first we need to cast

	if envsArr, ok := envsInterface.([]interface{}); ok {
		for i := range envsArr {
			keyAndVal := envsArr[i].(string)
			keyAndValArr := strings.Split(keyAndVal, "=")
			key := keyAndValArr[0]
			val := keyAndValArr[1]
			result[key] = val
		}
	} else if envsArr, ok := envsInterface.(map[string]interface{}); ok {
		// Key-value pair
		// Key: value
		for key, val := range envsArr {
			value, isString := val.(string)
			if !isString {
				value = strconv.Itoa(val.(int))
			}
			result[key] = value
		}
	}
	return result

}

func mapImage(image interface{}) string {
	if image == nil {
		return ""
	}
	return image.(string)
}

//Parse a docker compse file. Return the services
func parse(composeFilePath string) []Service {

	yaml, err := ioutil.ReadFile(composeFilePath)
	if err != nil {
		log.Fatal(err)
	}

	parsedYaml, err := loader.ParseYAML(yaml)
	if err != nil {
		log.Fatal(err)
	}
	var result []Service

	servicesUntyped := parsedYaml["services"].(map[string]interface{})
	for name, contentInterface := range servicesUntyped {
		content := contentInterface.(map[string]interface{})
		var service = Service{
			Name:        name,
			Image:       mapImage(content["image"]),
			Build:       mapBuild(content["build"]),
			Environment: mapEnviornment(content["environment"]),
			Volumes:     mapVolumes(content["volumes"]),
		}
		result = append(result, service)
	}

	return result
}
