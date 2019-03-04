package common

import (
	"crypto/md5"
	"fmt"
	"io"
)

func Md5(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Md5_blob(blob []byte) string {
	h := md5.New()
	h.Write(blob)
	return fmt.Sprintf("%x", h.Sum(nil))
}
