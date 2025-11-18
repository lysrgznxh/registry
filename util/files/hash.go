package files

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"io"
	"os"
)

func FileMd5(file string) (md5_ string, err_ error) {
	f, err_ := os.Open(file)
	if err_ != nil {
		return
	}

	buf := md5.New()
	if _, err_ = io.Copy(buf, f); err_ != nil {
		return
	}
	md5_ = hex.EncodeToString(buf.Sum(nil))
	return
}

func FileSha512(file string) (hash string, err error) {
	f, err_ := os.Open(file)
	if err_ != nil {
		return
	}

	buf := sha512.New()
	if _, err_ = io.Copy(buf, f); err_ != nil {
		return
	}
	hash = hex.EncodeToString(buf.Sum(nil))
	return
}
