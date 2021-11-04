package nio

import (
	"bytes"
	"io"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestWriteFileData(t *testing.T) {
	tests := []struct {
		name    string
		given FileData    
		wantContent   []byte
		wantErr bool
	}{
		{
			name: "not empty name and not empty content ",
			given: FileData{
				Name: "a.txt",
				Content: bytes.NewReader([]byte {1, 2, 3, 4, 5}), 
			},
			wantContent: []byte {5, 97, 46, 116, 120, 116, 1, 2, 3, 4, 5},
			wantErr: false,
		},
		{
			name: "empty content",
			given: FileData{
				Name: "a.txt",
				Content: bytes.NewReader([]byte {}), 
			},
			wantContent: []byte {5, 97, 46, 116, 120, 116},
			wantErr: false,
		},
		{
			name: "empty name",
			given: FileData{
				Name: "",
				Content: bytes.NewReader([]byte {1, 2, 3, 4}), 
			},
			wantContent: []byte {0, 1, 2, 3, 4},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if _, err := tt.given.WriteTo(w); (err != nil) != tt.wantErr {
				t.Errorf("WriteNamedStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotContent := w.Bytes(); !reflect.DeepEqual(gotContent, tt.wantContent) {
				t.Errorf("WriteNamedStream() = %v, want %v", gotContent, tt.wantContent)
			}
		})
	}
}

func TestReadFileData(t *testing.T) {
	type res struct {
		name string
		content []byte
	}
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    res
		wantErr bool
	}{
		{
			name: "not empty name and not empty content ",
			args: args {
				r: bytes.NewReader([]byte { 5, 97, 46, 116, 120, 116, 1, 2, 3, 4, 5}),
			},
			want: res {
				name: "a.txt",
				content: []byte {1, 2, 3, 4, 5},
			},
			wantErr: false,
		},
		{
			name: "empty name ",
			args: args {
				r: bytes.NewReader([]byte { 0, 1, 2, 3, 4, 5}),
			},
			want: res {
				name: "",
				content: []byte {1, 2, 3, 4, 5},
			},
			wantErr: false,
		},
		{
			name: "empty content",
			args: args {
				r: bytes.NewReader([]byte { 5, 97, 46, 116, 120, 116}),
			},
			want: res {
				name: "a.txt",
				content: []byte {},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadFileData(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadNamedStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotContent, _ := ioutil.ReadAll(got.Content)
			gotRes := res { name: got.Name, content: gotContent }
			if !reflect.DeepEqual(gotRes, tt.want) {
				t.Errorf("ReadNamedStream() = {name=%v,content=%v}, want {name=%v,content=%v}",
					gotRes.name, gotRes.content, tt.want.name, tt.want.content)
			}
		})
	}
}
