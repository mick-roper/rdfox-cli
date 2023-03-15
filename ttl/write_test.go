package ttl

import (
	"bytes"
	"testing"
)

func TestWrite(t *testing.T) {
	type args struct {
		data triples
	}
	tests := []struct {
		name    string
		args    args
		wantDst string
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				data: triples{
					"abc": map[string][]string{
						"xyz": {
							"one",
							"two",
							"three",
						},
					},
				},
			},
			wantDst: "abc\n\txyz\tone\t;\n\txyz\ttwo\t;\n\txyz\tthree\t.\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := &bytes.Buffer{}
			if err := Write(tt.args.data, dst); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDst := dst.String(); gotDst != tt.wantDst {
				t.Errorf("Write() = %v, want %v", gotDst, tt.wantDst)
			}
		})
	}
}
