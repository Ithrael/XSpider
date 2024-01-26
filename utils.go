package main

import (
	"regexp"
	"strings"
)

func IsRegexMatch(matchStr string, patternArr []string) bool {
	for i := 0; i < len(patternArr); i++ {
		match, _ := regexp.MatchString(patternArr[i], matchStr)
		if match {
			return true
		}
	}
	return false
}

func IsSubDomain(domain string, targets []string) bool {
	for i := 0; i < len(targets); i++ {
		domain = "." + strings.TrimLeft(domain, ".")
		target := "." + strings.TrimLeft(targets[i], ".")
		match := strings.HasSuffix(domain, target)
		if match {
			return true
		}
	}
	return false
}
