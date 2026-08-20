package pathinfo

import "testing"

func TestClassify(t *testing.T) {
	tests := []struct {
		path string
		want Kind
	}{
		{`C:\work\repo`, WindowsDrive},
		{`\\wsl.localhost\Ubuntu\home\me\repo`, UNC},
		{`/home/me/repo`, POSIX},
		{`/mnt/c/work/repo`, WSLMount},
		{`relative/path`, Unknown},
	}
	for _, test := range tests {
		if got := Classify(test.path); got != test.want {
			t.Errorf("Classify(%q) = %q, want %q", test.path, got, test.want)
		}
	}
}
