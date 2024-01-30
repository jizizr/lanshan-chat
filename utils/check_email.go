package utils

import "regexp"

var emailChecker = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+.[a-zA-Z]{2,}$`)

func Check_email(email string) bool {
	return emailChecker.MatchString(email)
}
