package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRunCmd(t *testing.T) {
	cmd := []string{"/bin/misterbin", "./testdata/echo.sh", "arg1=1", "arg2=2"}

	code := RunCmd(cmd, Environment{})

	require.Equal(t, 1, code)
}
