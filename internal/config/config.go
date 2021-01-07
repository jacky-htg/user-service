package config

import (
	"io/ioutil"
	"os"
	"strings"
)

//Setup environment from file .env
func Setup(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	datas := strings.Split(string(data), "\n")
	for _, env := range datas {
		e := strings.Split(env, "=")
		if len(e) >= 2 {
			os.Setenv(strings.TrimSpace(e[0]), strings.TrimSpace(strings.Join(e[1:], "=")))
		}
	}

	return nil
}
