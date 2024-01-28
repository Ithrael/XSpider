package main

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("test/test_config.yaml")
	if err != nil {
		t.Fatalf("Error loading config: %s", err)
	}

	// 验证Headers字段
	expectedUserAgent := "TestUserAgent"
	if config.Headers["user-agent"] != expectedUserAgent {
		t.Fatalf("Expected User-Agent to be %s, got %s", expectedUserAgent, config.Headers["user-agent"])
	}

	expectedAccept := "test/accept"
	if config.Headers["accept"] != expectedAccept {
		t.Fatalf("Expected Accept to be %s, got %s", expectedAccept, config.Headers["accept"])
	}

	expectedCookie := "test"
	if config.Headers["cookie"] != expectedCookie {
		t.Fatalf("Expected Cookie to be %s, got %s", expectedCookie, config.Headers["cookie"])
	}

	// 验证Restriction字段
	if config.Restriction.MaxDepth != 1 {
		t.Fatalf("Expected MaxDepth to be 0, got %d", config.Restriction.MaxDepth)
	}

	if config.Restriction.MaxCount != 1 {
		t.Fatalf("Expected MaxCount to be 0, got %d", config.Restriction.MaxCount)
	}

	// 验证AllowedDomains字段
	expectedAllowedDomains := []string{"example.com"}
	if !equalSlice(config.Restriction.AllowedDomains, expectedAllowedDomains) {
		t.Fatalf("Expected AllowedDomains to be %v, got %v", expectedAllowedDomains, config.Restriction.AllowedDomains)
	}
}

// 比较两个字符串切片是否相等
func equalSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
