package util

import (
	"encoding/json"
	"os"
)

type Server struct {
	Host	string	`json:"host"`
	Port	string	`json:"port"`
	Tunnel	string	`json:"tunnel"`
}

type Http struct {
	Path		string	`json:"path"`
	ServiceHost	string	`json:"serviceHost"`
	ServicePort	string	`json:"servicePort"`
}

type Conf struct {
	Server	Server	`json:"server"`
	Http	[]Http	`json:"http"`
}

// 读取yaml配置文件
// path string 文件路径
// v interface{} 格式化类型
func GetYaml(path string, v interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&v)
	return err
}