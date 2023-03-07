package auth

import (
	"os"
	"strings"
)

// UserParser Parses a string and return a userName/userHash. An error if the format of the string is incorrect.
type UserParser func(user string) (string, string, error)

const (
	authorizationHeader = "Authorization"
)

func getUsers(appendUsers []string, parser UserParser) (map[string]string, error) {
	users, err := loadUsers(appendUsers)
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]string)
	for _, user := range users {
		userName, userHash, err := parser(user)
		if err != nil {
			return nil, err
		}
		userMap[userName] = userHash
	}

	return userMap, nil
}

func loadUsers(appendUsers []string) ([]string, error) {
	var users []string

	return append(users, appendUsers...), nil
}

func getLinesFromFile(filename string) ([]string, error) {
	dat, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Trim lines and filter out blanks
	rawLines := strings.Split(string(dat), "\n")
	var filteredLines []string
	for _, rawLine := range rawLines {
		line := strings.TrimSpace(rawLine)
		if line != "" && !strings.HasPrefix(line, "#") {
			filteredLines = append(filteredLines, line)
		}
	}

	return filteredLines, nil
}
