package utils

import (
	"baka-go/constants"
	"baka-go/types"
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func GetAccessInfo() types.AccessInfo {
	f, err := os.Open(constants.DATA_DIRECTORY + "/" + constants.ACCESS_INFO_FILE)
	if err != nil {
		return types.AccessInfo{}
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var entries []string

	for scanner.Scan() {
		entries = append(entries, scanner.Text())
	}

	expiresIn, err := strconv.Atoi(entries[3])
	if err != nil {
			fmt.Println("Error reading access info.")
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
