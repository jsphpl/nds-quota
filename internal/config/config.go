package config

import (
	"os"
	"path"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	TemplateDirectory string `yaml:"templateDirectory"`
	DataDirectory     string `yaml:"dataDirectory"`
	Logfile           string `yaml:"logfile"`
}

func Load() (*Config, error) {
	dir := path.Dir(os.Args[0])
	f, err := os.Open(path.Join(dir, "config.yml"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var c Config
	err = yaml.NewDecoder(f).Decode(&c)
	return &c, err
}
