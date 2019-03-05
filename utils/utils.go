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

	"github.com/panshiqu/framework/define"
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
