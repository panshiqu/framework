package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

// Date 日期
func Date() int {
	y, m, d := time.Now().Date()
	return y*10000 + int(m)*100 + d
}

// ReadJSON 打开读取解析JSON文件
func ReadJSON(name string, js interface{}) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}

	defer f.Close()

	body, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, js)
}
