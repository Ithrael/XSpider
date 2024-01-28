package main

import (
	"reflect"
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
			name:       "test case 1 - Path match path with single pattern",
			matchStr:   "/path/test",
			patternArr: []string{"//path//.*", ".*path"},
			want:       true,
		},
		{
			name:       "test case 2 - Path match path with multiple patterns",
			matchStr:   "/testpath/test",
			patternArr: []string{"./path//.*", ".*path"},
			want:       true,
		},
		{
			name:       "test case 3 - Path match path with failure test ",
			matchStr:   "testpath/path1",
			patternArr: []string{".*//path//.*", ".*path$"},
			want:       false,
		},
		{
			name:       "test case 4 - Query string match queryKey with single pattern",
			matchStr:   "f=live",
			patternArr: []string{"src=.*typed_query", "f=live"},
			want:       true,
		},
		{
			name:       "test case 5 - Query string match queryKey with multiple patterns",
			matchStr:   "src=typed_query?f=live",
			patternArr: []string{"src=.*typed_query\\?f=live", "f=live"},
			want:       true,
		},
		{
			name:       "test case 6 - Query string match queryKey with failure test",
			matchStr:   "f=liv",
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

func TestReadUrlsFromFile(t *testing.T) {
	testFilePath := "./test/test_urls.txt"
	expected := []string{"https://www.apple.com", "https://support.apple.com", "https://music.apple.com"}

	result, err := ReadUrlsFromFile(testFilePath)
	if err != nil {
		t.Fatalf("Read File error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v，got %v", expected, result)
	}
}
