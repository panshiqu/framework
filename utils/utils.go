package utils

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	"../define"
	"github.com/tidwall/gjson"
	"github.com/mitchellh/mapstructure"
)

// Date 日期
func Date() int {
	y, m, _ := time.Now().Date()
	return y*100 + int(m)
}

// Signature 签名
func Signature(timestamp int64) string {
	s := fmt.Sprintf("%s%d", define.Token, timestamp)
	ss := strings.Split(s, "")
	sort.Strings(ss)
	sha := sha1.New()
	io.WriteString(sha, strings.Join(ss, ""))
	return fmt.Sprintf("%x", sha.Sum(nil))
}

// ReadJSON 打开读取解析JSON文件
func ReadJSON(name string, js interface{}) error {
	body,err := readFileContent(name)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, js)
}

// 读取文件内容
func readFileContent(path string)([]byte,error) {
	f, err := os.Open(path)
	if err != nil {
		return nil,err
	}

	defer f.Close()

	body, err := ioutil.ReadAll(f)
	if err != nil {
		return nil,err
	}
	return body,nil
}

func InitConfig(path string, config *define.GConfig) error {
	content, err := readFileContent(path)
	if err != nil {
		return err
	}
	m, _ := gjson.Parse(string(content)).Value().(map[string]interface{})

	//读取db信息
	mapstructure.Decode(m["db"], &config.DB)
	mapstructure.Decode(m["login"], &config.Login)
	mapstructure.Decode(m["game"], &config.Game)

	return nil
}
