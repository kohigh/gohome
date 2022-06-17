package main

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestReadNotExistingDir(t *testing.T) {
	_, err := ReadDir("./testdata/env/env")
	require.Error(t, err)
}

func TestReadInvalidFile(t *testing.T) {
	dir, err := ioutil.TempDir("./testdata", "test")
	if err != nil {
		log.Fatal(err)
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			log.Println("err")
		}
	}(dir)

	f, err := os.CreateTemp(dir, "*=.txt")
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Println("err")
		}
	}(f.Name())

	_, err = f.Write([]byte("asd"))
	if err != nil {
		log.Println(err)
	}

	env, err := ReadDir(dir)

	require.NoError(t, err)
	require.Zero(t, len(env))
}
