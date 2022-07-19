package util

import "testing"

func TestRemoveAudioFromWord(t *testing.T) {
	type args struct {
		word string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				word: "revise[sound:revise]",
			},
			want: "revise",
		},
		{
			args: args{
				word: "revise",
			},
			want: "revise",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveAudioFromWord(tt.args.word); got != tt.want {
				t.Errorf("RemoveAudioFromWord() = %v, want %v", got, tt.want)
			}
		})
	}
}
