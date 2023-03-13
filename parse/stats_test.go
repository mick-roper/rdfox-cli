package parse

// import (
// 	"io"
// 	"reflect"
// 	"strings"
// 	"testing"
// )

// func TestStats(t *testing.T) {
// 	const happyPath = `1	"Component name"	"hello"
// 1	"world"	1
// 2	"component name"	"rick"
// 2	"and"	"morty"`

// 	type args struct {
// 		r io.Reader
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want statistics
// 	}{
// 		{
// 			name: "happy path",
// 			args: args{
// 				r: strings.NewReader(""),
// 			},
// 			want: statistics{
// 				"hello": map[string]interface{}{
// 					"world": 1,
// 				},
// 				"rick": map[string]interface{}{
// 					"and": "morty",
// 				},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := Stats(tt.args.r); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Stats() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
