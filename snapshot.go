package xattrdb

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	xattr "github.com/pkg/xattr"
)

var snapshots []string = make([]string, 0)

func GetSnapshots() []string {
	return snapshots
}

func CreateSnapshot() string {
	now := time.Now()
	nsec := strconv.FormatInt(now.UnixNano(), 10)
	var sb strings.Builder
	sb.WriteString(GetPath())
	sb.WriteString(nsec)
	os.OpenFile(sb.String(), os.O_RDONLY|os.O_CREATE, 0666)
	snapshots = append(snapshots, nsec)
	sb.Reset()
	return nsec
}

func ReadSnapshot(snapshot, key string) (string, error) {
	data, err := xattr.Get(path+snapshot, GetName(key))
	if err != nil {
		return ReadData(key)
	}
	value, err := Decompress(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return value, nil
}

func writeSnapshot(snapshot, key, value string) bool {
	compressed, err := Compress(value)
	if err != nil {
		log.Println(err)
		return false
	}
	if err = xattr.Set(path+snapshot, GetName(key), compressed); err != nil {
		log.Println(err)
		return false
	}
	return true
}

func CopyOnWrite(key string) error {
	if _, err := ReadData(key); err != nil {
		return nil
	}
	for _, snapshot := range snapshots {
		_, err := xattr.Get(path+snapshot, GetName(key))
		if err != nil {
			value, _ := ReadData(key)
			if !writeSnapshot(snapshot, key, value) {
				return errors.New("Snapshot Copy on Write failure")
			}
		}
	}
	return nil
}
