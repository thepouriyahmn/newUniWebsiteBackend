package auth

import "regexp"

type Regex struct {
}

func NewRegex() Regex {
	return Regex{}
}

func (r Regex) IsValidPassword(password string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9]{6,}$`)
	return re.MatchString(password)
}
