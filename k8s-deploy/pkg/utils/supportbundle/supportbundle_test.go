package supportbundle

import (
	"os"
	"reflect"
	"testing"
)

const (
	bundirDir = "test-bundle"
)

var (
	file = File{"1-2-3", []byte("body")}
)

func TestNewFile(t *testing.T) {
	type args struct {
		name string
		body []byte
	}
	tests := []struct {
		name string
		args args
		want *File
	}{
		{"newFile", args{"1 2 3", []byte("body")}, &file},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFile(tt.args.name, tt.args.body); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pathExists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"pathExists", args{"."}, true, false},
		{"pathNotExists", args{bundirDir}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pathExists(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("pathExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("pathExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckDir(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"checkExsitedDir", args{"."}, false},
		{"checkDirNotExist", args{bundirDir}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exist, _ := pathExists(tt.args.dir)
			if err := CheckDir(tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("CheckDir() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !exist {
				os.Remove(tt.args.dir)
			}
		})
	}
}

func Test_saveFile(t *testing.T) {
	type args struct {
		file *File
		dir  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"saveFile", args{&file, bundirDir}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckDir(tt.args.dir); err != nil {
				t.Errorf("Test_saveFile CheckDir() error = %v", err)
			}
			if err := saveFile(tt.args.file, tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("saveFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := os.RemoveAll(tt.args.dir); err != nil {
				t.Errorf("Test_saveFile RemoveAll() error = %v", err)
			}
		})
	}
}

func Test_dirToZip(t *testing.T) {
	type args struct {
		dir string
		dst string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"dirToZip", args{bundirDir, "test.zip"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckDir(tt.args.dir); err != nil {
				t.Errorf("Test_saveFile CheckDir() error = %v", err)
			}
			if err := dirToZip(tt.args.dir, tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("dirToZip() error = %v, wantErr %v", err, tt.wantErr)
			}
			if exist, err := pathExists(tt.args.dst); !exist || err != nil {
				t.Errorf("dirToZip() pathExists error = %v", err)
			}
			if err := os.RemoveAll(tt.args.dir); err != nil {
				t.Errorf("Test_dirToZip RemoveAll() dir error = %v", err)
			}
			if err := os.RemoveAll(tt.args.dst); err != nil {
				t.Errorf("Test_dirToZip RemoveAll() dst error = %v", err)
			}
		})
	}
}
