package util

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Serverhost string	`json:"serverhost"`
	Serverport int	`json:"serverport"`
	Tunnel string	`json:"name"`
	Rules []Rule	`json:"rules"`
}

type Rule struct {
	Tag string	`json:"tag"`
	Host string	`json:"host"`
	Port int	`json:"port"`
}

// 读取yaml配置文件
// path string 文件路径
// v interface{} 格式化类型
func (t *Config) GetConfig(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("yamlFile.Get error #%v\n", err)
		return nil, err
	}
	err = yaml.Unmarshal(file, t)
	if err != nil {
		fmt.Printf("Unmarshal: %v", err)
		return nil, err
	}
	return t, err
}