package base

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
)

func GetJsKV(src string, key string) string {
	var val string
	r, err := regexp.Compile(fmt.Sprintf(`%v:"([^"]+)"`, key))
	if err != nil {
		log.Printf("Regex did not compile when trying to get key %v from js. %v", key, err)
	}
	if m := r.FindStringSubmatch(src); m != nil {
		val = m[1]
	}

	return val
}

func findJsArray(src string) (string, error) { //TODO write test
	n := 0
	for i, c := range src {
		if c == '[' {
			n++
		}
		if c == ']' {
			if n < 1 {
				break
			}
			n--
		}
		if n == 0 {
			return src[:i+1], nil
		}
	}

	return "", errors.New("cannot find array: invalid js")
}

func GetJsArrayVar(src string, name string) (string, error) {
	var val string
	r, err := regexp.Compile(fmt.Sprintf(`%v\s*=\s*`, name))
	if err != nil {
		return val, err
	}
	if m := r.FindStringIndex(src); m != nil {
		return findJsArray(src[m[1]:])
	}

	return val, fmt.Errorf("array variable %v not found", name)
}

func FetchFileContent(url string) (string, error) {
	var content string
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to get js file: %v", err)
		return content, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return content, err
	}
	content = string(body)
	if res.StatusCode != 200 {
		return content, fmt.Errorf("failed to get %v. response: %v", url, content)
	}

	return content, nil
}
