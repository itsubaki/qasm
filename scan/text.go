package scan

import (
	"bufio"
	"fmt"
	"io"
)

func Text(r io.Reader) (string, error) {
	scanner := bufio.NewScanner(r)

	var text string
	for scanner.Scan() {
		text += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("scan: %w", err)
	}

	return text, nil
}

func MustText(r io.Reader) string {
	return Must(Text(r))
}

func Must(s string, err error) string {
	if err != nil {
		panic(err)
	}

	return s
}
