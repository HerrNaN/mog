package mog

import "testing"

func Test_lockFilePathOf(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "same directory",
			args: args{"file.txt"},
			want: ".#file.txt",
		},
		{
			name: "different directory",
			args: args{"some/path/file.txt"},
			want: "some/path/.#file.txt",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lockFilePathOf(tt.args.filePath); got != tt.want {
				t.Errorf("lockFilePathOf() = %v, want %v", got, tt.want)
			}
		})
	}
}
