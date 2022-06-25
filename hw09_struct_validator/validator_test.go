package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	invalidErrTests := []struct {
		in           interface{}
		expectedErrs []error
	}{
		{
			in: User{
				ID:     "123",
				Age:    0,
				Email:  "vasya",
				Role:   "customer",
				Phones: []string{"123", "234"},
			},
			expectedErrs: []error{invalidLen, invalidMin, invalidDueToRegexp, invalidIn},
		},
		{
			in:           App{Version: "123"},
			expectedErrs: []error{invalidLen},
		},
		{
			in:           Response{Code: 300, Body: "asd"},
			expectedErrs: []error{invalidIn},
		},
	}

	for i, tt := range invalidErrTests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			e := Validate(tt.in)
			for _, err := range tt.expectedErrs {
				require.ErrorContains(t, e, err.Error())
			}
		})
	}

	validErrTests := []struct {
		in           interface{}
		expectedErrs []error
	}{
		{
			in: User{
				ID:     "98739125-daf0-42a8-a58a-e9a29ba2d75b",
				Age:    20,
				Email:  "vasya@gmail.com",
				Role:   "admin",
				Phones: []string{"12312312312", "12312312311", "12312312313"},
			},
		},
		{
			in: App{Version: "1.232"},
		},
		{
			in: Token{Header: []byte{}, Payload: []byte{}, Signature: []byte{}},
		},
		{
			in: Response{Code: 404, Body: "asd"},
		},
	}

	for i, tt := range validErrTests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			e := Validate(tt.in)
			require.Nil(t, e)
		})
	}
}
