package blog

import (
	"hack/cmd/util"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateWebpFormat(t *testing.T) {
	defer os.Remove("./testdate/imageforwebp/simple.webp")
	GenerateWebpFormat("./testdate/imageforwebp/simple.png")
}

func TestIsImage(t *testing.T) {
	type args struct {
		fPath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "true:png",
			args: args{fPath: "./testdate/sourcedir/blog1/simple.png"},
			want: true,
		},
		{
			name: "true:webp",
			args: args{fPath: "./testdate/sourcedir/blog1/simple.webp"},
			want: true,
		},
		{
			name: "false:text",
			args: args{fPath: "./testdate/sourcedir/blog1/textfile.txt"},
			want: false,
		},
		{
			name: "false:textWithPNGSuffix",
			args: args{fPath: "./testdate/sourcedir/blog1/textfile.png"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := IsImage(tt.args.fPath); got != tt.want {
				t.Errorf("IsImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCopyImgFromDirOfBlogToDirImgCDN(t *testing.T) {
	is := assert.New(t)
	defer os.RemoveAll("./testdate/targetdir/blog1")
	o := &SyncImageOption{
		CompressEnable: true,
	}
	o.CopyImgFromDirOfBlogToDirImgCDN("./testdate/sourcedir", "./testdate/targetdir")

	_, err := os.Stat(filepath.Join("./testdate/targetdir", "blog1/simple.png"))
	is.NoError(err)

	_, err = os.Stat(filepath.Join("./testdate/targetdir", "blog1/simple.webp"))
	is.NoError(err)
}

func TestCompressImageByCommandTool(t *testing.T) {
	is := assert.New(t)
	image, err := NewImage("./testdate/compress/simple.png", PNG)
	is.NoError(err)

	err = CompressImageByCommandTool(*image)
	is.NoError(err)
}

func TestSyncImageOption_Run(t *testing.T) {
	cfg := util.BlogConfig{
		BlogDir:       "",
		BlogSourceDir: "./testdate/sourcedir",
		ImageRepoDir:  "./testdate/targetdir",
	}
	o := SyncImageOption{
		CompressEnable: true,
		ReleaseEnable:  false,
		Update:         false,
	}
	o.Run(cfg)
}