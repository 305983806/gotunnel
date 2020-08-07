package main

import (
	"encoding/json"
	"os"
)

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
	err = decoder.Decode(v)
	return err
}