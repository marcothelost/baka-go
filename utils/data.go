package utils

import (
	"baka-go/constants"
	"baka-go/types"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func getAccessInfoEntries(f *os.File) []string {
	scanner := bufio.NewScanner(f)
	var entries []string

	for scanner.Scan() {
		entries = append(entries, scanner.Text())
	}

	return entries
}

func GetEndpointURL(route string) string {
	return constants.BAKALARI_ENDPOINT + route
}

func GetAccessInfo() types.AccessInfo {
	f, err := os.Open(constants.DATA_DIRECTORY + "/" + constants.ACCESS_INFO_FILE)
	if err != nil {
		return types.AccessInfo{}
	}
	defer f.Close()

	var entries []string = getAccessInfoEntries(f)
	expiresIn, err := strconv.Atoi(entries[3])
	if err != nil {
		fmt.Println("Error reading access info!")
		return types.AccessInfo{}
	}

	return types.AccessInfo{
		Access_Token: entries[0],
		Refresh_Token: entries[1],
		Token_Type: entries[2],
		Expires_In: expiresIn,
	}
}

func SaveAccessInfo(acessInfo types.AccessInfo) {
	os.Mkdir(constants.DATA_DIRECTORY, os.ModePerm)

	f, err := os.Create(constants.DATA_DIRECTORY + "/" + constants.ACCESS_INFO_FILE)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString(acessInfo.Access_Token + "\n")
	f.WriteString(acessInfo.Refresh_Token + "\n")
	f.WriteString(acessInfo.Token_Type + "\n")
	f.WriteString(fmt.Sprint(acessInfo.Expires_In) + "\n")
}

func HandleExpiredToken() types.AccessInfo {
	f, err := os.Open(constants.DATA_DIRECTORY + "/" + constants.ACCESS_INFO_FILE)
	if err != nil {
		fmt.Println("Pro použití BakaGO se musíte přihlásit!")
		fmt.Println("-- bakago login")
		os.Exit(1)
	}
	defer f.Close()

	var entries []string = getAccessInfoEntries(f)
	var refreshToken string = entries[1]

	postData := []byte(
		fmt.Sprintf("client_id=ANDR&grant_type=refresh_token&refresh_token=%v", refreshToken),
	)
	buffer := bytes.NewBuffer(postData)

	res, err := http.Post(GetEndpointURL(constants.LOGIN_ROUTE), "application/x-www-form-urlencoded", buffer)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println("Vypršel váš refresh token! Musíte se znou přihlásit!")
		fmt.Println("-- bakago login")
		return types.AccessInfo{}
	}

	content, _ := io.ReadAll(res.Body)

	var responseData types.AccessInfo
	_ = json.Unmarshal(content, &responseData)
	SaveAccessInfo(responseData)

	return responseData
}
