package collectors

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
)

func GetStringFromJs(src string, key string) string {
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
