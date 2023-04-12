package logger

import (
	"testing"
)

func TestNew(t *testing.T) {
	for _, test := range []struct {
		dest      string
		wantError bool
	}{
		{
			dest:      "stderr",
			wantError: false,
		},
		{
			dest:      "stdout",
			wantError: false,
		},
		{
			dest:      "/tmp/gtpl.log",
			wantError: false,
		},
		{
			dest:      "/non/existing/path/to/file",
			wantError: true,
		},
	} {
		_, err := New(test.dest)
		gotErr := err != nil
		if gotErr != test.wantError {
			if err != nil {
				t.Errorf("New(%q) = %v, got error: %v, want error: %v", test.dest, err, gotErr, test.wantError)
			}
		}
	}
}
