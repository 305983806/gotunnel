package util

import (
	"fmt"
	"github.com/305983806/gotunnel/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"testing"
)

func TestConfig_GetConfig(t *testing.T) {
	path := "./config.yaml"
	outputFile(path)

	var config = util.Config{}
	config.GetConfig(path)
	fmt.Printf("%#v", config)

	deleteFile(path)
}

func outputFile(path string) {
	var config = util.Config{
		ServerHost: "127.0.0.1",
		ServerPort: 8002,
		Tunnel:     "cp",
		Rules:      []util.Rule{
			{
				Tag:  "neo",
				Host: "127.0.0.1",
				Port: 8080,
			},
			{
				Tag: "authentication",
				Host: "127.0.0.1",
				Port: 8081,
			},
		},
	}
	c, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(path, c, 0777)
	if err != nil {
		fmt.Println("write file error: ", err)
	}
}

func deleteFile(path string) {
	os.Remove(path)
}