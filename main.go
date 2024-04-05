package main

import (
	"baka-go/types"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	BAKALARI_ENDPOINT = "https://sosro.bakalari.cz"
)

const (
	LOGIN_COMMAND = "login"
	TIMETABLE_COMMAND = "timetable"
	AVERAGES_COMMAND = "averages"
)

const (
	LOGIN_ROUTE = "/api/login"
)

func getEndpointURL(route string) string {
	return BAKALARI_ENDPOINT + route
}

func login() {
	var username string
	var password string

	fmt.Print("Login\n\n")
	fmt.Print("Username: ")
	fmt.Scan(&username)
	fmt.Print("Password: ")
	fmt.Scan(&password)
	fmt.Println()

	postData := []byte(
		fmt.Sprintf(`client_id=ANDR&grant_type=password&username=%v&password=%v`, username, password),
	)
	buffer := bytes.NewBuffer(postData)

	response, err := http.Post(getEndpointURL(LOGIN_ROUTE), "application/x-www-form-urlencoded", buffer)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	content, _ := io.ReadAll(response.Body)

	if (response.StatusCode != 200) {
		var data types.BakalariErrorResponse
		_ = json.Unmarshal(content, &data)
		
		if (data.Error == "invalid_grant") {
			println("Neplatné uživatelské jméno nebo heslo!")
			os.Exit(1)
		}
	}

	var data types.BakalariLoginResponse
	_ = json.Unmarshal(content, &data)
	fmt.Println(data.Access_Token)
}

func timetable() {
	fmt.Println("Timetable")
}

func averages() {
	fmt.Println("Averages")
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

	switch (args[0]) {
	case LOGIN_COMMAND:
		login()
	case TIMETABLE_COMMAND:
		timetable()
	case AVERAGES_COMMAND:
		averages()
	default:
		fmt.Println("Command not found")
		os.Exit(1)
	}
}
