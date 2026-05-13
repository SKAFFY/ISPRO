package main

import (
	"fmt"
	"regexp"

	"github.com/dlclark/regexp2"
)

var (
	simplePattern   = regexp.MustCompile(`^[^:]+:[^:]+:[^:]+$`)
	complexPassword = regexp2.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*])[a-zA-Z\d!@#$%^&*]{8,}$`, 0)
)

func ValidateSimpleLink(link string) bool {
	return simplePattern.MatchString(link)
}

func ValidateComplexPassword(password string) bool {
	match, _ := complexPassword.MatchString(password)
	return match
}

func main() {
	links := []string{
		"source:target:description",
		"invalid",
		"a:b:c:d",
		"",
	}

	fmt.Println("=== Simple Pattern (source:target:description) ===")
	for _, link := range links {
		fmt.Printf("%q -> %v\n", link, ValidateSimpleLink(link))
	}

	fmt.Println("\n=== Complex Pattern (Password with lookahead validation) ===")
	passwords := []string{
		"Password1!",
		"weak",
		"noNumbers!",
		"NoSpecial1",
		"Valid1!Pass",
		"",
	}
	for _, pwd := range passwords {
		fmt.Printf("%q -> %v\n", pwd, ValidateComplexPassword(pwd))
	}
}