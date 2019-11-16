package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/aureleoules/gocaml/models"
	"github.com/bwmarrin/discordgo"
)

func removeLastLine(str string) string {
	arr := strings.Split(str, "\n")
	return strings.Join(arr[:len(arr)-2], "\n")
}

// IsCodeEvaluation checks if discord message contains a code evaluation request, and return sanitized code
func IsCodeEvaluation(m *discordgo.MessageCreate) (bool, string, string) {
	reg := regexp.MustCompile("(?s)```(ocaml|python)(.*?)```")
	match := reg.FindStringSubmatch(m.Content)

	if match == nil {
		return false, "", ""
	}

	lang := match[1]

	code := match[len(match)-1]
	code = strings.Replace(code, "\"", "\\"+"\"", -1)
	return true, code, lang
}

// FormatEvaluation formats CAML evaluation
func FormatEvaluation(eval string) string {
	formatted := strings.Replace(eval, "        ", "", -1)
	formatted = strings.Replace(formatted, "   ", " ", -1)
	formatted = removeLastLine(formatted)
	return formatted
}

// ContainsError checks if evaluation contains an error
func ContainsError(eval string) bool {
	lines := strings.Split(eval, "\n")
	return strings.Contains(lines[len(lines)-1], "Error: ")
}

// IsStats cehcks if discord message is a stat request, and return user id if specified
func IsStats(m *discordgo.MessageCreate) (bool, string) {
	if strings.Contains(m.Content, prefix) {
		reg := regexp.MustCompile("<@!([0-9]*?)>")
		match := reg.FindStringSubmatch(m.Content)
		if match != nil {
			return true, match[1]
		}
		return true, ""
	}
	return false, ""
}

// ParseStats parses stats
func ParseStats(users []models.User) string {
	result := "```\n"
	for _, user := range users {
		result += user.Username + "#" + user.Discriminator + ":\n"
		result += "Success: " + strconv.Itoa(user.SuccessCount) + "\n"
		result += "Errors: " + strconv.Itoa(user.ErrorCount) + "\n"
		result += "Last evaluation: " + user.LastEvaluation.String() + "\n\n"
	}
	result += "```"
	return result
}

func VerifyPythonCode(code string) error {
	illegalModules := []string{
		"os",
	}
	illegalFunctions := []string{
		"open",
		"exec",
		"__import__",
	}

	for _, m := range illegalModules {
		r, _ := regexp.MatchString(`^[^#"]?\s*((import|from)([\s,_\w\\])*`+m+`)`, code)
		if r {
			return errors.New("Illegal module imported: '" + m + "'")
		}
	}

	for _, f := range illegalFunctions {
		r, _ := regexp.MatchString(`^[^#"]*`+f+`\s*\(.*`, code)
		if r {
			return errors.New("Illegal function used: '" + f + "'")
		}
	}

	return nil
}
