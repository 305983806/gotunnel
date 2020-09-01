package util

import (
	"fmt"
	"testing"
)

func TestConfig_GetConfig(t *testing.T) {
	var config = Config{}
	config.GetConfig("../config.yaml")
	fmt.Printf("%v", config)
}