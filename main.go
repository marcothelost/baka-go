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

func help() {
	blue := color.New(color.FgBlue).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Printf("%v %v - Spustí proces přihlášení\n", blue("bakago"), green("login"))
	fmt.Printf("%v %v - Vypíše průměry známek\n", blue("bakago"), green("marks"))
	fmt.Printf("  ○ marks %v - Za výpis přidá seznam známek\n", yellow("-l"))
}

func marks(flags []string) {
	var marks = utils.FetchData[types.MarksListing](constants.MARKS_ROUTE)

	var maxLen int = 0
	for _, subject := range marks.Subjects {
		if utf8.RuneCountInString(subject.Subject.Name) > maxLen {
			maxLen = utf8.RuneCountInString(subject.Subject.Name)
		}
	}

	blue := color.New(color.FgBlue).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	for _, subject := range marks.Subjects {
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

func final(flags []string) {
	var finalMarks = utils.FetchData[types.FinalMarksListing](constants.FINAL_ROUTE)
	uniqueSubjects := make(map[string]*types.Subject)

	for _, term := range finalMarks.CertificateTerms {
		for _, subject := range term.Subjects {
			uniqueSubjects[subject.Id] = &subject
		}	
	}

	var maxLen int = 0
	for _, subject := range uniqueSubjects {
		if utf8.RuneCountInString(subject.Abbrev) > maxLen {
			maxLen = utf8.RuneCountInString(subject.Abbrev)
		}
	}

	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	blue.Printf("%-*s ", maxLen + 2, "")
	for _, term := range finalMarks.CertificateTerms {
		fmt.Printf("%v ", green(term.Grade))
	}
	fmt.Println()


	for _, subject := range uniqueSubjects {
		blue.Printf("%-*s", maxLen, subject.Abbrev)
		fmt.Print(" - ")
		for _, term := range finalMarks.CertificateTerms {
			var printedMark string = " "
			for _, mark := range term.FinalMarks {
				if mark.SubjectId == subject.Id {
					if len(mark.MarkText) == 0 {
						printedMark = "?"
						continue
					}
					printedMark = mark.MarkText
				}
			}
			fmt.Printf("%v ", printedMark)
		}
		if !slices.Contains(flags, "l") {
			fmt.Println()
			continue
		}
		fmt.Printf("- %v", yellow(subject.Name))
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
		color.Blue("BakaGo " + constants.PROJECT_VERSION)
		fmt.Println("-- bakago help")
		return
	}

	if (args[0] == constants.LOGIN_COMMAND) {
		login()
		return
	}

	switch (args[0]) {
	case constants.HELP_COMMAND:
		help()
	case constants.MARKS_COMMAND:
		marks(flags)
	case constants.FINAL_COMMAND:
		final(flags)
	default:
		color.Red("Tento příkaz neexistuje!")
		color.Red("-- bakago help")
		os.Exit(1)
	}
}
