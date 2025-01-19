package utils

import "testing"

func TestIsLikelyBinary(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{"ascii text", []byte("This is a normal ASCII text file with printable characters"), false},
		{"text with newlines", []byte("Line 1\nLine 2\nLine 3"), false},
		{"mixed content", []byte("Hello\x00World"), false},
		{"highly binary", []byte{0x00, 0xFF, 0x01, 0x02, 0xFE, 0xFF}, true},
		{"mostly binary", []byte{0x7F, 0x45, 0x00, 0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFF, 'a'}, true},
		{"single char", []byte{'a'}, false},
		{"single binary", []byte{0x00}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsLikelyBinary(tt.data)
			if got != tt.expected {
				t.Errorf("IsLikelyBinary() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestShouldIncludeFile(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		exclude       string
		include       string
		startDir      string
		shouldInclude bool
	}{
		{"hidden file", ".git", "", "", ".", false},
		{"included pattern", "test.go", "", "*.go", ".", true},
		{"excluded pattern", "test.txt", "*.txt", "", ".", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShouldIncludeFile(tt.path, tt.exclude, tt.include, tt.startDir)
			if got != tt.shouldInclude {
				t.Errorf("ShouldIncludeFile() = %v, want %v for test %s", got, tt.shouldInclude, tt.name)
			}
		})
	}

	moreTests := []struct {
		name          string
		path          string
		exclude       string
		include       string
		startDir      string
		shouldInclude bool
	}{
		{"both patterns", "test.go", "*.txt", "*.go", ".", true},
		{"no patterns", "normal.file", "", "", ".", true},
		{"exclude wins", "test.txt", "*.txt", "*.txt", ".", false},
		{"different dir", "subdir/test.go", "", "*.go", ".", true},
		{"hidden directory", ".git/config", "", "", ".", false},
		{"multiple extensions", "test.go.txt", "*.txt", "*.go", ".", false},
	}

	for _, tt := range moreTests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShouldIncludeFile(tt.path, tt.exclude, tt.include, tt.startDir)
			if got != tt.shouldInclude {
				t.Errorf("ShouldIncludeFile() = %v, want %v for test %s", got, tt.shouldInclude, tt.name)
			}
		})
	}
}
