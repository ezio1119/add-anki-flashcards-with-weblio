package weblio

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestQuery(t *testing.T) {
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
				Level:       1,
				AudioURL:    "https://weblio.hs.llnwd.net/e7/img/dict/kenej/audio/S-C1393A8_E-C13AF80.mp3",
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
