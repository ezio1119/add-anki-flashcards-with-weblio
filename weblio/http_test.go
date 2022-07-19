package weblio

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestQueryEver(t *testing.T) {
	type args struct {
		ctx   context.Context
		query string
	}
	tests := []struct {
		name    string
		args    args
		want    *QueryResult
		wantErr bool
	}{
		{
			name: "ever",
			args: args{
				ctx:   context.TODO(),
				query: "ever",
			},
			want: &QueryResult{
				Query:       "ever",
				Description: "かつて、いつか、これまでに、かつて(…することがない)、決して(…ない)、いずれ、これまで、今まで、いつも、常に",
				Examples: examples{
					&example{En: "Have you ever seen a tiger?", Ja: "トラを見たことがありますか 《★この文の答えは Yes, I have (once). または No, I have not. または No, I never have.》."},
					&example{En: "Did you ever see him while you were in Tokyo?", Ja: "東京にいる間に彼に会いましたか."},
					&example{En: "How can I ever thank you (enough)?", Ja: "ほんとにお礼の申しようもありません."},
					&example{En: "Few tourists ever come to this part of the country.", Ja: "この地方にまで来る観光旅行者はほとんどない."},
					&example{En: "I won't ever forget you.", Ja: "決して君のことは忘れない."},
					&example{En: "They won't ever forget it.", Ja: "彼らは決してそれを忘れることはないよ."},
					&example{En: "If you (should) ever come this way, be sure to call on us.", Ja: "もしこちらへおいでになることがおありでしたら, 必ず私たちの所に寄ってください."},
					&example{En: "If I ever catch him!", Ja: "彼を捕らえようものなら(ただではおかないぞ)!"},
					&example{En: "He is [was] a great musician if ever there was one.", Ja: "彼こそ正に大音楽家だ[だった]."},
					&example{En: "It's raining harder than ever.", Ja: "雨がさらにいっそう激しく降っている."},
					&example{En: "He's the greatest poet that England ever produced.", Ja: "彼は英国が生んだ最も偉大な詩人だ."},
					&example{En: "He's ever quick to respond.", Ja: "彼はいつも応答が早い."},
					&example{En: "ever‐active", Ja: "常に活動的な."},
					&example{En: "ever‐present", Ja: "常に存在する."},
					&example{En: "What ever is she doing?", Ja: "彼女は一体何をしているのか."},
					&example{En: "Who ever can it be?", Ja: "一体だれだろう."},
					&example{En: "Why ever did you not say so?", Ja: "一体なぜそう言わなかったのだ."},
					&example{En: "Is this ever beautiful!", Ja: "これは実に美しいではないか."},
				},
				Level:    1,
				AudioURL: "https://weblio.hs.llnwd.net/e7/img/dict/kenej/audio/S-C1393A8_E-C13AF80.mp3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Query(tt.args.ctx, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Query() mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}

func Test_examples_String(t *testing.T) {
	tests := []struct {
		name string
		exx  examples
		want string
	}{
		{
			name: "ever",
			exx: examples{
				&example{En: "Have you ever seen a tiger?", Ja: "トラを見たことがありますか 《★この文の答えは Yes, I have (once). または No, I have not. または No, I never have.》."},
				&example{En: "Did you ever see him while you were in Tokyo?", Ja: "東京にいる間に彼に会いましたか."},
			},
			want: "Have you ever seen a tiger?: トラを見たことがありますか 《★この文の答えは Yes, I have (once). または No, I have not. または No, I never have.》.<br>Did you ever see him while you were in Tokyo?: 東京にいる間に彼に会いましたか.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.exx.String(); got != tt.want {
				t.Errorf("examples.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
