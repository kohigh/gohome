package main

import (
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	fromFile, err := os.CreateTemp("./", "*.txt")
	defer func() {
		os.Remove(fromFile.Name())
		os.Remove("1.txt")
	}()

	if err != nil {
		log.Println(err)
	}

	tests := []struct {
		name   string
		from   string
		to     string
		limit  int64
		offset int64
	}{
		{name: "Negative offset", from: fromFile.Name(), to: "1.txt", offset: -1},
		{name: "Wrong fromFile", from: "asd", to: "1.txt"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := Copy(tc.from, tc.to, tc.offset, tc.limit)
			require.Error(t, err)
		})
	}
}
