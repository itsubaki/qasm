package io

import (
	"bufio"
	"io"
)

func Scan(r io.Reader) (string, error) {
	scanner := bufio.NewScanner(r)

	var text string
	for scanner.Scan() {
		text += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return text, nil
}

func MustScan(r io.Reader) string {
	return Must(Scan(r))
}

func Must(s string, err error) string {
	if err != nil {
		panic(err)
	}

	return s
}
