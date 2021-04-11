package server

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
		assert.Equal(t, true, CreateData("foo", "bar"))
	})
	t.Run("reading a string", func(t *testing.T) {
		actual, err := ReadData("foo")
		assert.Nil(t, err)
		assert.Equal(t, "bar", actual)
	})
	t.Run("updating a string", func(t *testing.T) {
		assert.Equal(t, true, UpdateData("foo", "qaz"))
	})
	t.Run("reading a string after update", func(t *testing.T) {
		actual, err := ReadData("foo")
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

func TestSnapshot(t *testing.T) {
	t.Run("storing a string", func(t *testing.T) {
		assert.Equal(t, true, CreateData("foo", "bar"))
	})
	t.Run("reading a string after update", func(t *testing.T) {
		actual, err := ReadData("foo")
		assert.Nil(t, err)
		assert.Equal(t, "bar", actual)
	})
	firstSnapshot := CreateSnapshot()
	t.Run("storing a string", func(t *testing.T) {
		assert.Equal(t, true, UpdateData("foo", "qaz"))
	})
	t.Run("reading a string after update", func(t *testing.T) {
		actual, err := ReadData("foo")
		assert.Nil(t, err)
		assert.Equal(t, "qaz", actual)
	})
	secondSnapshot := CreateSnapshot()
	t.Run("reading a string from a snapshot", func(t *testing.T) {
		actual, err := ReadSnapshot(firstSnapshot, "foo")
		assert.Nil(t, err)
		assert.Equal(t, "bar", actual)
	})
	DeleteData("foo")
	t.Run("reading a string from a snapshot", func(t *testing.T) {
		actual, err := ReadSnapshot(firstSnapshot, "foo")
		assert.Nil(t, err)
		assert.Equal(t, "bar", actual)
	})
	t.Run("reading a string from a snapshot", func(t *testing.T) {
		actual, err := ReadSnapshot(secondSnapshot, "foo")
		assert.Nil(t, err)
		assert.Equal(t, "qaz", actual)
	})
	CreateData("qaz", "zaq")
	t.Run("reading a string not held in a snapshot", func(t *testing.T) {
		actual, _ := ReadSnapshot(firstSnapshot, "qaz")
		assert.Equal(t, "", actual)
	})
}
