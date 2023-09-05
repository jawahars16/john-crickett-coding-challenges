package conf

import (
	"io"
	"os"

	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v2"
)

type Conf struct {
	Port     int      `yaml:"port"`
	Backends []string `yaml:"backends"`
	Health   Health   `yaml:"health"`
}

type Health struct {
	URL      string `yaml:"url"`
	Interval int    `yaml:"interval"`
}

func Load(file string) (*Conf, error) {
	confFile, err := os.Open("conf.yml")
	defer confFile.Close()
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	var conf Conf
	confData, err := io.ReadAll(confFile)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	yaml.Unmarshal(confData, &conf)
	return &conf, nil
}
