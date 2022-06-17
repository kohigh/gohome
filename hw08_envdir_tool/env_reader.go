package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func (v Environment) getEnvPairs() []string {
	var envs []string
	for key, val := range v {
		if val.NeedRemove {
			err := os.Unsetenv(key)
			if err != nil {
				log.Printf("Env %v can't be unset due to the error: %v", key, err)
			}
		}

		envs = append(envs, fmt.Sprintf("%v=%v", key, val.Value))
	}

	return append(os.Environ(), envs...)
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	envs := make(Environment)

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range entries {
		if strings.Contains(f.Name(), "=") {
			continue
		}

		file, err := os.Open(fmt.Sprintf("%v/%v", dir, f.Name()))
		if err != nil {
			log.Printf("File %v opened with error: %v\n", f.Name(), err)
			continue
		}
		buf := bufio.NewReader(file)

		value, err := buf.ReadBytes('\n')
		if err != nil && err != io.EOF {
			log.Printf("File %v read with error: %v\n", f.Name(), err)
			continue
		}
		value = bytes.TrimRight(value, "\n")
		value = bytes.TrimRight(value, "\t")
		value = bytes.TrimRight(value, " ")
		value = bytes.Replace(value, []byte("\x00"), []byte("\n"), -1)

		var env = EnvValue{NeedRemove: true}

		if len(value) != 0 {
			env.Value = string(value)
			env.NeedRemove = false
		}

		envs[f.Name()] = env
	}

	return envs, nil
}
