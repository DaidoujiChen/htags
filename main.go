package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"google.golang.org/appengine"
)

const (
	rawURL = "https://raw.githubusercontent.com/wiki/Mapaler/EhTagTranslator/database/"
	index  = "rows"
)

func isEmpty(s string) bool {
	return len(strings.Replace(s, " ", "", -1)) == 0
}

func clearString(s string) string {
	preAndSufSpaceReg := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
	insideSpaceReg := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	removeImageReg := regexp.MustCompile(`(?im)!\[.*]\(.*\)`)
	result := preAndSufSpaceReg.ReplaceAllString(s, "")
	result = insideSpaceReg.ReplaceAllString(result, " ")
	result = removeImageReg.ReplaceAllString(result, "")
	return result
}

func fetchCategories(r *http.Request) ([]string, error) {
	categories := []string{}

	// 取得 category 列表網址們
	res, err := httpGet(rawURL+index+".md", r)
	if err != nil {
		return categories, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return categories, err
	}

	// 取出需要的網址段落
	urlReg := regexp.MustCompile(`(?im)https:\/\/github\.com\/Mapaler\/EhTagTranslator\/wiki\/[a-zA-Z]*`)
	urls := urlReg.FindAllString(string(body), -1)
	for _, url := range urls {
		splits := strings.Split(url, "/")
		lastComponet := splits[len(splits)-1]
		category := strings.Replace(lastComponet, ".md", "", -1)
		categories = append(categories, category)
	}
	return categories, nil
}

func fetchContentsIn(category string, r *http.Request) (map[string]string, error) {
	contents := map[string]string{}

	// 取得該 category 的 raw data
	res, err := httpGet(rawURL+category+".md", r)
	if err != nil {
		return contents, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return contents, err
	}

	// 切出 table 的部分
	clearContentReg := regexp.MustCompile(`(?im)^\s*(\|.*\|)\s*$`)
	lines := clearContentReg.FindAllString(string(body), -1)
	for _, line := range lines {
		splits := strings.Split(line, "|")
		if len(splits) != 6 {
			continue
		}

		key := clearString(splits[1])
		value := clearString(splits[2])
		if !isEmpty(key) && !isEmpty(value) {
			contents[key] = value
		}
	}
	return contents, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := fetchCategories(r)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	allContents := map[string]string{}
	for _, category := range categories {
		contents, err := fetchContentsIn(category, r)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}

		for k, v := range contents {
			if exist, ok := allContents[k]; ok {
				if strings.Contains(exist, v) {
					continue
				}

				allContents[k] = exist + "|" + v
				continue
			}
			allContents[k] = v
		}
	}

	jsonString, err := json.Marshal(allContents)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

func main() {
	http.HandleFunc("/", indexHandler)
	// http.ListenAndServe(":8030", nil)
	appengine.Main()
}
