package openai

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

type OpenAIConfig struct {
	OpenAIAPIKey string `yaml:"OPENAI_API_KEY"`
}

func parseConfig() (*OpenAIConfig, error) {
	confFile := "config.yaml"
	path, err := getPath()
	if err != nil {
		return nil, err
	}

	file := strings.Join([]string{path, confFile}, string(os.PathSeparator))

	if !fsExists(path) {
		if err := os.Mkdir(path, 0777); err != nil {
			return nil, err
		}
	}

	if !fsExists(file) {
		getOpenAIKey(file)
	}

	f, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var config OpenAIConfig
	if err := yaml.Unmarshal(f, &config); err != nil {
		log.Fatal(err)
	}

	return &config, nil
}

func fsExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func getOpenAIKey(filePath string) {
	var prompt string

	fmt.Println("Paste in your chatGPT API Token. To obtain one:")
	fmt.Println("https://platform.openai.com/docs/api-reference/authentication")

	_, err := fmt.Scanln(&prompt)
	if err != nil {
		log.Fatal(err)
	}

	config, err := yaml.Marshal(&OpenAIConfig{prompt})
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(filePath, config, 0666); err != nil {
		log.Fatal(err)
	}
}

func getPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", nil
	}

	var paths []string = make([]string, 2)

	if runtime.GOOS == "windows" {
		paths = []string{home, "%PROGRAMDATA%", "aicli"}
	} else {
		paths = []string{home, ".config", "aicli"}
	}

	path := strings.Join(paths, string(os.PathSeparator))

	return path, nil
}
