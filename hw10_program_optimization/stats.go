package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)

	scanner := bufio.NewScanner(r)

	emailDelim := "@"
	u := User{}

	for scanner.Scan() {
		text := scanner.Text()

		if strings.Contains(text, domain) {
			err := u.UnmarshalJSON(scanner.Bytes())
			if err != nil {
				continue
			}

			if !strings.Contains(u.Email, domain) {
				continue
			}

			from := strings.Index(u.Email, emailDelim) + 1
			key := strings.ToLower(u.Email[from:])
			result[key]++
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
