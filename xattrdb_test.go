package xattrdb

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	SetPath("/tmp/xattrdb")
	SetShards(2)
	CreateShards()
}

func TestStorage(t *testing.T) {
	t.Run("storing a string", func(t *testing.T) {
		assert.Equal(t, true, DataUpdate("foo", "bar"))
	})
	t.Run("reading a string", func(t *testing.T) {
		actual, err := DataRead("foo")
		assert.Nil(t, err)
		assert.Equal(t, "bar", actual)
	})
	t.Run("updating a string", func(t *testing.T) {
		assert.Equal(t, true, DataUpdate("foo", "qaz"))
	})
	t.Run("reading a string after update", func(t *testing.T) {
		actual, err := DataRead("foo")
		assert.Nil(t, err)
		assert.Equal(t, "qaz", actual)
	})
}

func TestSharding(t *testing.T) {
	t.Run("sharding a key", func(t *testing.T) {
		assert.Equal(t, GetPath()+"0", Shard("qaz"))
	})
	t.Run("sharding a key", func(t *testing.T) {
		assert.Equal(t, GetPath()+"1", Shard("foo"))
	})
}
