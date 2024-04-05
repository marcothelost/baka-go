package main

import (
	"baka-go/constants"
	"baka-go/types"
	"baka-go/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"unicode/utf8"
)

func getEndpointURL(route string) string {
	return constants.BAKALARI_ENDPOINT + route
}

func login() {
	var username string
	var password string

	fmt.Println("Přihlášení do BakaGo")
	fmt.Println("Vaše údaje nejsou zasílány žádné třetí straně.")
	fmt.Println()
	fmt.Print("Uživatelské jméno: ")
	fmt.Scan(&username)
	fmt.Print("Heslo: ")
	fmt.Scan(&password)
	fmt.Println()

	postData := []byte(
		fmt.Sprintf(`client_id=ANDR&grant_type=password&username=%v&password=%v`, username, password),
	)
	buffer := bytes.NewBuffer(postData)

	res, err := http.Post(getEndpointURL(constants.LOGIN_ROUTE), "application/x-www-form-urlencoded", buffer)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	content, _ := io.ReadAll(res.Body)

	if (res.StatusCode != 200) {
		var responseData types.BakalariError
		_ = json.Unmarshal(content, &responseData)
		
		if (responseData.Error == "invalid_grant") {
			println("Neplatné uživatelské jméno nebo heslo!")
			os.Exit(1)
		}
	}

	var responseData types.AccessInfo
	_ = json.Unmarshal(content, &responseData)
	utils.SaveAccessInfo(responseData)

	fmt.Println("Byli jste úspěšně přihlášeni.")
}

func marks() {
	var accessInfo types.AccessInfo = utils.GetAccessInfo()

	req, err := http.NewRequest("GET", getEndpointURL(constants.MARKS_ROUTE), nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer " + accessInfo.Access_Token)

	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	content, _ := io.ReadAll(res.Body)

	var responseData types.MarksListing
	_ = json.Unmarshal(content, &responseData)

	var maxLen int = 0
	for _, subject := range responseData.Subjects {
		if utf8.RuneCountInString(subject.Subject.Name) > maxLen {
			maxLen = utf8.RuneCountInString(subject.Subject.Name)
		}
	}

	for _, subject := range responseData.Subjects {
		var paddedName string = fmt.Sprintf("%-*s", maxLen, subject.Subject.Name)
		fmt.Printf("%v - %v - %v\n", subject.Subject.Abbrev, paddedName, subject.AverageText)
	}
}

func main() {
	var args []string
	var flags []string

	for _, arg := range os.Args[1:] {
		if len(arg) > 1 && arg[0] == '-' {
			flags = append(flags, arg[1:])
			continue
		}
		args = append(args, arg)
	}

	if len(args) == 0 {
		fmt.Println("No arguments!")
		os.Exit(1)
	}

	if (args[0] == constants.LOGIN_COMMAND) {
		login()
		return
	}

	switch (args[0]) {
	case constants.MARKS_COMMAND:
		marks()
	default:
		fmt.Println("Command not found")
		os.Exit(1)
	}
}
