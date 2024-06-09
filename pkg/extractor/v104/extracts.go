package v104

import (
	"errors"
	"regexp"
)

func extractUpgradeToken(pageHTML []byte) (string, error) {
	return ExtractToken(pageHTML)
}

func extractTearDownToken(pageHTML []byte) (string, error) {
	return ExtractToken(pageHTML)
}

func ExtractToken(pageHTML []byte) (string, error) {
	rgx := regexp.MustCompile(`var token = "([^"]+)"`)
	m := rgx.FindSubmatch(pageHTML)
	if len(m) != 2 {
		return "", errors.New("unable to find token")
	}
	return string(m[1]), nil
}
