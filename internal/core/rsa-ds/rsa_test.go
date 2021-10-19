package rsa_ds

import (
	. "math/big"
	"reflect"
	"testing"
)

var constHash = func(value *Int) HashFn {
	return func(orig *Int) (*Int, error) {
		return value, nil
	}
}

func TestSign(t *testing.T) {
	tests := []struct {
		name    string
		key     *SecretKey
		msg     *Int
		hashFn  HashFn
		want    *Signed
		wantErr bool
	}{
		{
			name:    "1",
			key:     &SecretKey{N: NewInt(11 * 29), Exp: NewInt(11)},
			msg:     NewInt(4),
			want:    &Signed{Msg: NewInt(4), Signature: NewInt(92)},
			hashFn:  constHash(NewInt(4)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sign(tt.key, tt.msg, tt.hashFn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sign() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSignatureValid(t *testing.T) {
	tests := []struct {
		name    string
		key     *PubKey
		signed  *Signed
		hashFn  HashFn
		want    bool
		wantErr bool
	}{
		{
			name:    "1",
			key:     &PubKey{N: NewInt(11 * 29), Exp: NewInt(51)},
			signed:  &Signed{Msg: NewInt(4), Signature: NewInt(92)},
			hashFn:  constHash(NewInt(4)),
			want:    true,
			wantErr: false,
		},
		{
			name:    "2",
			key:     &PubKey{N: NewInt(11 * 29), Exp: NewInt(51)},
			signed:  &Signed{Msg: NewInt(4), Signature: NewInt(97)},
			hashFn:  constHash(NewInt(4)),
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsSignatureValid(tt.key, tt.signed, tt.hashFn)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsSignatureValid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsSignatureValid() got = %v, want %v", got, tt.want)
			}
		})
	}
}
