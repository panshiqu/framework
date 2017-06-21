package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// Date 日期
func Date() string {
	y, m, d := time.Now().Date()
	return fmt.Sprintf("%d%02d%02d", y, m, d)
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
