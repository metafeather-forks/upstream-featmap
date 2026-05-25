package main

import (
	"testing"
)

func TestWorkspaceNameIsValid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"valid short", "ab", true},
		{"valid normal", "myworkspace", true},
		{"valid with numbers", "workspace123", true},
		{"too short", "a", false},
		{"too long", string(make([]byte, 201)), false},
		{"reserved account", "account", false},
		{"reserved link", "link", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := workspaceNameIsValid(tt.input); got != tt.want {
				t.Errorf("workspaceNameIsValid(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateTitle(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantErr   bool
		wantTitle string
	}{
		{"valid", "My Project", false, "My Project"},
		{"trimmed", "  hello  ", false, "hello"},
		{"too short", "", true, ""},
		{"too long", string(make([]byte, 201)), true, ""},
		{"exactly right", "A", false, "A"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateTitle(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateTitle(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.wantTitle {
				t.Errorf("validateTitle(%q) = %q, want %q", tt.input, got, tt.wantTitle)
			}
		})
	}
}

func TestColorIsValid(t *testing.T) {
	valid := []string{"WHITE", "GREY", "RED", "ORANGE", "YELLOW", "GREEN", "TEAL", "BLUE", "INDIGO", "PURPLE", "PINK"}
	for _, c := range valid {
		if !colorIsValid(c) {
			t.Errorf("colorIsValid(%q) should be true", c)
		}
	}
	if colorIsValid("INVALID") {
		t.Error("colorIsValid(INVALID) should be false")
	}
	if colorIsValid("") {
		t.Error("colorIsValid(\"\") should be false")
	}
}

func TestAreAnnotationsValid(t *testing.T) {
	if !areAnnotationsValid("") {
		t.Error("empty annotations should be valid")
	}
	if !areAnnotationsValid("RISKY") {
		t.Error("RISKY should be valid")
	}
	if !areAnnotationsValid("RISKY,BLOCKED,IDEA") {
		t.Error("multiple valid annotations should be valid")
	}
	if areAnnotationsValid("INVALID") {
		t.Error("invalid annotation should not be valid")
	}
	if areAnnotationsValid("RISKY,INVALID") {
		t.Error("partially invalid annotations should not be valid")
	}
}

func TestLevelIsValid(t *testing.T) {
	valid := []string{"VIEWER", "EDITOR", "ADMIN", "OWNER"}
	for _, l := range valid {
		if !levelIsValid(l) {
			t.Errorf("levelIsValid(%q) should be true", l)
		}
	}
	if levelIsValid("SUPERUSER") {
		t.Error("levelIsValid(SUPERUSER) should be false")
	}
	if levelIsValid("") {
		t.Error("levelIsValid(\"\") should be false")
	}
}

func TestIsEditor(t *testing.T) {
	if isEditor("VIEWER") {
		t.Error("VIEWER should not be editor")
	}
	if !isEditor("EDITOR") {
		t.Error("EDITOR should be editor")
	}
	if !isEditor("ADMIN") {
		t.Error("ADMIN should be editor")
	}
	if !isEditor("OWNER") {
		t.Error("OWNER should be editor")
	}
}
