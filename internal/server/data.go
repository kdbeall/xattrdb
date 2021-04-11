package server

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"

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

func GetName(key string) string {
	return prefix + key
}

func CreateShards() {
	var sb strings.Builder
	for i := 0; int64(i) < shards; i++ {
		sb.WriteString(GetPath())
		locationNum := strconv.Itoa(i)
		sb.WriteString(locationNum)
		os.OpenFile(sb.String(), os.O_RDONLY|os.O_CREATE, 0666)
		sb.Reset()
	}
}

func CreateData(key, value string) bool {
	_, err := xattr.Get(Shard(key), GetName(key))
	if err == nil {
		return false
	}
	compressed, err := Compress(value)
	if err != nil {
		log.Println(err)
		return false
	}
	CopyOnWriteNil(key)
	if err = xattr.Set(Shard(key), GetName(key), compressed); err != nil {
		log.Println(err)
		return false
	}
	return true
}

func ReadData(key string) (string, error) {
	data, err := xattr.Get(Shard(key), GetName(key))
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

func ReadCompressed(key string) ([]byte, error) {
	data, err := xattr.Get(Shard(key), GetName(key))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return data, nil
}

func UpdateData(key, value string) bool {
	_, err := xattr.Get(Shard(key), GetName(key))
	if err != nil {
		log.Println(err)
		return false
	}
	compressed, err := Compress(value)
	if err != nil {
		log.Println(err)
		return false
	}
	CopyOnWrite(key)
	if err = xattr.Set(Shard(key), GetName(key), compressed); err != nil {
		log.Println(err)
		return false
	}
	return true
}

func DeleteData(key string) bool {
	CopyOnWrite(key)
	if err := xattr.Remove(Shard(key), GetName(key)); err != nil {
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
