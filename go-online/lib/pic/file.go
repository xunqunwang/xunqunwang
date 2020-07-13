package pic

import (
	"encoding/base64"
	// "io/ioutil"
	"os"
)

func Base64Encoding(name string) (cc string, err error) {
	var (
		ff *os.File
		n  int
	)
	if ff, err = os.Open(name); err != nil {
		return
	}
	defer ff.Close()
	src := make([]byte, 500000)
	if n, err = ff.Read(src); err != nil {
		return
	}
	//base64压缩
	cc = base64.StdEncoding.EncodeToString(src[:n])
	return
}

func Base64Decoding(s string) (dist []byte, err error) {
	//解压
	if dist, err = base64.StdEncoding.DecodeString(string(s)); err != nil {
		return
	}
	return
}
