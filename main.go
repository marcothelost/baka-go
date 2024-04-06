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
	"slices"
	"unicode/utf8"

	"github.com/fatih/color"
)

func login() {
	var username string
	var password string

	blue := color.New(color.FgBlue).SprintFunc()

	fmt.Printf("Přihlášení do %s\n", blue("BakaGo"))
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

	res, err := http.Post(utils.GetEndpointURL(constants.LOGIN_ROUTE), "application/x-www-form-urlencoded", buffer)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	content, _ := io.ReadAll(res.Body)

	if (res.StatusCode != 200) {
		var responseData types.BakalariError
		_ = json.Unmarshal(content, &responseData)
		
		if (responseData.Error == "invalid_grant") {
			color.Red("Neplatné uživatelské jméno nebo heslo!")
			os.Exit(1)
		}
	}

	var responseData types.AccessInfo
	_ = json.Unmarshal(content, &responseData)
	utils.SaveAccessInfo(responseData)

	color.Green("Byli jste úspěšně přihlášeni.")
}

func marks(flags []string) {
	var accessInfo types.AccessInfo = utils.GetAccessInfo()

	req, err := http.NewRequest("GET", utils.GetEndpointURL(constants.MARKS_ROUTE), nil)
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

	if (res.StatusCode == 401) {
		accessInfo = utils.HandleExpiredToken()
		req.Header.Set("Authorization", "Bearer " + accessInfo.Access_Token)

		res, err = client.Do(req)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
	}

	content, _ := io.ReadAll(res.Body)

	var responseData types.MarksListing
	_ = json.Unmarshal(content, &responseData)

	var maxLen int = 0
	for _, subject := range responseData.Subjects {
		if utf8.RuneCountInString(subject.Subject.Name) > maxLen {
			maxLen = utf8.RuneCountInString(subject.Subject.Name)
		}
	}

	blue := color.New(color.FgBlue).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	for _, subject := range responseData.Subjects {
		var paddedName string = fmt.Sprintf("%-*s", maxLen, subject.Subject.Name)
		fmt.Printf("%v - %v - %v", blue(subject.Subject.Abbrev), paddedName, green(subject.AverageText))
		if !slices.Contains(flags, "l") {
			fmt.Println()
			continue
		}
		fmt.Print("- ")
		for _, mark := range subject.Marks {
			fmt.Printf("%v ", yellow(mark.MarkText))
		}
		fmt.Println()
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
		marks(flags)
	default:
		color.Red("Tento příkaz neexistuje!")
		color.Red("-- bakago help")
		os.Exit(1)
	}
}
