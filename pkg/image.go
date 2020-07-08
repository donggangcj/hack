package pkg

import (
	"errors"
	"hack/cmd/util"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
)

type ImageType string

const (
	WindowsIcon   ImageType = "image/x-icon"
	WindowsCursor ImageType = "image/x-icon"
	BMPImage      ImageType = "image/bmp"
	GIF           ImageType = "image/gif"
	WEBP          ImageType = "image/webp"
	PNG           ImageType = "image/png"
	JPEG          ImageType = "image/jpeg"
)

var ValidateImageType = []ImageType{
	WindowsIcon, WindowsCursor, BMPImage, GIF, WEBP, PNG, JPEG,
}

// IsImage determine whether a file is a image
// https://stackoverflow.com/questions/25959386/how-to-check-if-a-file-is-a-valid-image
func IsImage(fPath string) (bool, ImageType) {
	open, err := os.Open(fPath)
	if err != nil {
		return false, ""
	}
	defer open.Close()

	buffer := make([]byte, 512)
	_, err = open.Read(buffer)
	if err != nil {
		return false, ""
	}
	contentType := http.DetectContentType(buffer)
	for _, imageType := range ValidateImageType {
		if contentType == string(imageType) {
			return true, imageType
		}
	}
	return false, ""
}

type Image struct {
	Path        string
	ImageFormat ImageType
	Blog        string
	Size        int64
}

//Name return base name of origin image
func (i *Image) Name() string {
	return strings.TrimSuffix(filepath.Base(i.Path), filepath.Ext(i.Path))
}

func (i *Image) PathWithWEBPExt() string {
	return strings.TrimSuffix(i.Path, filepath.Ext(i.Path)) + ".webp"
}

func (i *Image) NameWithExt() string {
	return filepath.Base(i.Path)
}

func (i *Image) NameWithBlog() string {
	return filepath.Join(i.Blog, i.NameWithExt())
}

func (i *Image) ToTableColumnString() []string {
	return []string{i.Blog, i.NameWithExt(), util.ByteCountSI(i.Size)}
}

func NewImage(fPath string, imageType ImageType) (*Image, error) {
	open, err := os.Stat(fPath)
	if err != nil {
		return nil, err
	}

	dirs := strings.Split(filepath.Dir(fPath), "/")
	blogName := dirs[len(dirs)-1]
	srcImg := &Image{
		Path:        fPath,
		ImageFormat: imageType,
		Blog:        blogName,
		Size:        open.Size(),
	}
	return srcImg, nil
}

type Images []Image

func (i Images) Print(writer io.Writer) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Blog", "Image", "Size"})
	table.SetRowLine(true)
	table.SetAutoMergeCells(true)
	//table.SetColumnColor(tablewriter.Colors{tablewriter.BgGreenColor}, tablewriter.Colors{tablewriter.BgBlueColor},
	//	tablewriter.Colors{tablewriter.BgHiRedColor})
	for _, image := range i {
		table.Append(image.ToTableColumnString())
	}
	//table.SetFooter([]string{"", "Total", strconv.Itoa(len(i))})
	table.Render()
}


//CompressImageByCommandTool compress image by the local command tool.
func CompressImageByCommandTool(image Image) error {
	var tool string
	switch image.ImageFormat {
	case PNG:
		tool = "optipng"
	case JPEG:
		tool = "jpegtopnm"
	default:
		tool = ""
	}
	if tool == "" {
		logrus.Errorf("not support image type:%s", image.ImageFormat)
		return errors.New("nonsupport image type")
	}
	cmd := exec.Command(tool, image.Path)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
