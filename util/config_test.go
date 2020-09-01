package util

import (
	"fmt"
	"github.com/305983806/gotunnel/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"testing"
)

var config = util.Config{
	Serverhost: "127.0.0.1",
	Serverport: 8002,
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

func TestConfig_GetConfig(t *testing.T) {
	path := "./config.yaml"
	outputFile(path)

	var conf = util.Config{}
	conf.GetConfig(path)

	str1, _ := yaml.Marshal(conf)
	str2, _ := yaml.Marshal(config)
	if string(str1) != string(str2) {
		t.Errorf("读取config.yaml结果未达预期！\n")
	}

	deleteFile(path)
}

func outputFile(path string) {
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