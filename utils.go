package main

import (
	"regexp"
)

func IsMatch(matchStr string, patternArr []string) bool {
	for i := 0; i < len(patternArr); i++ {
		match, _ := regexp.MatchString(patternArr[i], matchStr)
		if match {
			return true
		}
	}
	return false
}
