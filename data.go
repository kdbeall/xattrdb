package xattrdb

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strconv"

	xattr "github.com/pkg/xattr"
)

const (
	prefix = "user."
)

var shards int64
var path string

func SetShards(num int) {
	shards = int64(num)
}

func SetPath(fp string) {
	path = fp
}

func GetPath() string {
	return path
}

func CreateShards() {
	for i := 0; int64(i) < shards; i++ {
		locationNum := strconv.Itoa(i)
		os.OpenFile(path+locationNum, os.O_RDONLY|os.O_CREATE, 0666)
	}
}

func DataCreate(key, value string) bool {
	return DataUpdate(key, value)
}

func DataRead(key string) (string, error) {
	data, err := xattr.Get(Shard(key), prefix+key)
	if err != nil {
		log.Println(err)
		return "", err
	}
	value, err := Decompress(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return value, nil
}

func DataUpdate(key, value string) bool {
	compressed, err := Compress(value)
	if err != nil {
		log.Println(err)
		return false
	}
	if err = xattr.Set(Shard(key), prefix+key, compressed); err != nil {
		log.Println(err)
	}
	return true
}

func DataDelete(key string) bool {
	if err := xattr.Remove(Shard(key), prefix+key); err != nil {
		log.Println(err)
		return false
	}
	return true
}

func Compress(value string) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(value)); err != nil {
		log.Println(err)
		return nil, err
	}
	if err := gz.Close(); err != nil {
		log.Println(err)
		return nil, err

	}
	return b.Bytes(), nil
}

func Decompress(compressedValue []byte) (string, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(compressedValue))
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer gz.Close()
	data, err := ioutil.ReadAll(gz)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(data), nil
}

func Shard(key string) string {
	hKey := Hash(key)
	result, i, locations := new(big.Int), new(big.Int), new(big.Int)
	i = i.SetBytes(hKey)
	locations = locations.SetInt64(int64(shards))
	result.Mod(i, locations)
	return path + result.String()
}

func Hash(key string) []byte {
	h := sha1.New()
	h.Write([]byte(key))
	hKey := h.Sum(nil)
	return hKey
}
