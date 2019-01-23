package internal

import (
	"strings"
)

// ex. $MY_IMAGE=registry/imagename:version
type imageProperty struct {
	key      string
	registry string
	image    string
	version  string
}

type envVal struct {
	value string
}

// A docker compose service
type Service struct {
	Name        string
	Image       string
	Environment map[string]string
	Build       map[string]interface{}
	Volumes     []volum
}

type volum struct {
	localPath  string
	monutPath  string
	permission string
}

func (service *Service) getUsedSystemEnvs() []string {
	var result []string
	if service.Image != "" {
		imgKey := strings.TrimPrefix(service.Image, "$")
		result = append(result, imgKey)
	}

	for _, value := range service.Environment {
		if strings.ContainsRune(value, '$') { //Only allow single value. TODO: Extend to support ${VAL}
			result = append(result, strings.TrimPrefix(value, "$"))
		}
	}

	return result
}

func (service *Service) getImageEnv(allEnvs map[string]envVal, registry string) imageProperty {
	imgKey := strings.TrimPrefix(service.Image, "$")
	value := allEnvs[imgKey].value
	entry := imageProperty{
		key: imgKey,
	}

	if len(registry) > 0 {
		if registryAndName := strings.SplitAfter(value, registry); len(registryAndName) > 1 {
			entry.registry = registryAndName[0]
		}
	}

	withoutRegistry := strings.TrimPrefix(value, registry)
	nameAndVersion := strings.Split(withoutRegistry, ":")
	entry.image = nameAndVersion[0]
	if len(nameAndVersion) > 1 {
		entry.version = nameAndVersion[1]
	}

	return entry
}

func (img *imageProperty) getValue() string {
	result := img.registry + img.image
	if img.version != "" {
		result += ":" + img.version
	}
	return result
}

func (img *imageProperty) getKeyAndValue() string {
	return img.key + "=" + img.getValue()
}
