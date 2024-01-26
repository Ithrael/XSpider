package main

import (
	"testing"
)

func TestIsRegexMatch(t *testing.T) {
	tests := []struct {
		name       string
		matchStr   string
		patternArr []string
		want       bool
	}{
		{
			name:       "test case 1 - Domain match test with single pattern",
			matchStr:   "apple.com",
			patternArr: []string{"apple.com.cn", "apple.com"},
			want:       true,
		},
		{
			name:       "test case 2 - Domain match test with multiple patterns",
			matchStr:   "test.apple.com.cn",
			patternArr: []string{"*.apple.com.cn", "*.apple.com"},
			want:       true,
		},
		{
			name:       "test case 3 - Domain match failure test",
			matchStr:   "testapple.com.cn",
			patternArr: []string{"*.apple.com.cn", "*.apple.com"},
			want:       false,
		},
		{
			name:       "test case 4 - Path match path with single pattern",
			matchStr:   "/path/test",
			patternArr: []string{"/path/*", "*path"},
			want:       true,
		},
		{
			name:       "test case 5 - Path match path with multiple patterns",
			matchStr:   "/testpath/test",
			patternArr: []string{"/path/*", "*path*"},
			want:       true,
		},
		{
			name:       "test case 6 - Path match path with failure test ",
			matchStr:   "testpath/path1",
			patternArr: []string{"/path/*", "*path"},
			want:       false,
		},
		{
			name:       "test case 7 - Query string match queryKey with single pattern",
			matchStr:   "f=live",
			patternArr: []string{"src=*typed_query", "f=live"},
			want:       true,
		},
		{
			name:       "test case 8 - Query string match queryKey with multiple patterns",
			matchStr:   "src=typed_query",
			patternArr: []string{"src=*typed_query", "f=live"},
			want:       true,
		},
		{
			name:       "test case 9 - Query string match queryKey with failure test",
			matchStr:   "nf=live",
			patternArr: []string{"src=*typed_query", "f=live"},
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRegexMatch(tt.matchStr, tt.patternArr); got != tt.want {
				t.Errorf("IsRegexMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSubDomain(t *testing.T) {
	tests := []struct {
		name       string
		matchStr   string
		patternArr []string
		want       bool
	}{
		{
			name:       "test case 1 - Domain match test with single pattern",
			matchStr:   "apple.com",
			patternArr: []string{"apple.com.cn", "apple.com"},
			want:       true,
		},
		{
			name:       "test case 2 - Domain match test with multiple patterns",
			matchStr:   "test.apple.com.cn",
			patternArr: []string{"apple.com.cn", "apple.com"},
			want:       true,
		},
		{
			name:       "test case 3 - Domain match failure test",
			matchStr:   "testapple.com.cn",
			patternArr: []string{"apple.com.cn", "apple.com"},
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSubDomain(tt.matchStr, tt.patternArr); got != tt.want {
				t.Errorf("IsSubDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
