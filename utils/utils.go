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

	"github.com/panshiqu/framework/define"
)

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
		return Wrap(err)
	}

	defer f.Close()

	body, err := ioutil.ReadAll(f)
	if err != nil {
		return Wrap(err)
	}

	return Wrap(json.Unmarshal(body, js))
}

// CheckError 检查错误
func CheckError(data []byte) error {
	me := &define.MyError{}

	if err := json.Unmarshal(data, me); err != nil {
		return Wrap(err)
	}

	if me.Errno != define.ErrnoSuccess {
		return Wrap(me)
	}

	return nil
}
