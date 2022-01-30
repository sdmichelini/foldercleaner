package rules

import (
	"testing"
	"time"
)

func Test_folderNameForTime(t *testing.T) {

	tests := []struct {
		name string
		args time.Time
		want string
	}{
		{"Test Basic Folder", time.Date(2020, 2, 1, 1, 1, 1, 1, time.UTC), "2020/02"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := folderNameForTime(tt.args); got != tt.want {
				t.Errorf("folderNameForTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
