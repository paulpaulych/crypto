package key_gen

import (
	"github.com/paulpaulych/crypto/cmd/cli"
	"math/big"
	"reflect"
	"testing"
)

func TestConf_NewCmd(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    cli.Cmd
		wantErr cli.CmdConfError
	}{
		{
			name:    "1",
			args:    []string{"-p", "30803", "-q", "1297"},
			want:    &Cmd{p: big.NewInt(30803), q: big.NewInt(1297)},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &Conf{}
			got, err := conf.NewCmd(tt.args)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("error = %v, want %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cmd = %v, want %v", got, tt.want)
			}
		})
	}
}
