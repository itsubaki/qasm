package scan

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func Text(r io.Reader) (string, error) {
	scanner := bufio.NewScanner(r)

	var text strings.Builder
	for scanner.Scan() {
		text.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("scan: %w", err)
	}

	return text.String(), nil
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
