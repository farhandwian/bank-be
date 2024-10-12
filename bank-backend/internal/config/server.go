package config

import "fmt"

type serverConfig struct {
	Port         int  `yaml:"port" json:"port"`
	ReadTimeout  uint `yaml:"read_timeout" json:"read_timeout"`
	WriteTimeout uint `yaml:"write_timeout" json:"write_timeout"`
}

func (l serverConfig) Addr() string {
	return fmt.Sprintf(":%d", l.Port)
}
