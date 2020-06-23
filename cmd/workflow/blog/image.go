package blog

import (
	"errors"
	"fmt"
	"hack/cmd/util"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//MaxImageSize present the max size that a image need to be compress
const MaxImageSize = 2048000

type SyncImageOption struct {
	CompressEnable bool
	ReleaseEnable  bool
	// if update is true,the sync command will update images which existed
	Update bool
}

func NewSyncImageOption() *SyncImageOption {
	return &SyncImageOption{}
}

func NewSyncImage(cfg util.BlogConfig) *cobra.Command {
	o := NewSyncImageOption()
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "sync blog's image to local repo ,then release it",
		Run: func(cmd *cobra.Command, args []string) {
			o.Run(cfg)
		},
	}

	cmd.Flags().BoolVar(&o.CompressEnable, "compress", true, "compress image whose size is over 2MB")
	cmd.Flags().BoolVar(&o.ReleaseEnable, "release", true, "release images from local repo to remote")
	cmd.Flags().BoolVar(&o.Update, "force-update", false, "update existing images")
	return cmd
}

func (o *SyncImageOption) Run(cfg util.BlogConfig) {

	err := o.CopyImgFromDirOfBlogToDirImgCDN(cfg.BlogSourceDir, cfg.ImageRepoDir)
	PrintErrorToStdErr(err,"err happen when sync image")

	if o.ReleaseEnable {
		logrus.Info("🚀: release image")
		ReleaseImageTOCDN(cfg.ImageRepoRootDir)
	}
}

type Image struct {
	content     *os.File
	imageFormat ImageType
	blog        string
}

//Name return base name of origin image
func (i *Image) Name() string {
	return strings.TrimSuffix(filepath.Base(i.content.Name()), filepath.Ext(i.content.Name()))
}

func (i *Image) Path() string {
	return i.content.Name()
}

func (i *Image) PathWithWEBPExt() string {
	return strings.TrimSuffix(i.content.Name(), filepath.Ext(i.content.Name())) + ".webp"
}

func (i *Image) NameWithExt() string {
	return filepath.Base(i.content.Name())
}

func (i *Image) NameWithBlog() string {
	return filepath.Join(i.blog, i.NameWithExt())
}

func NewImage(fPath string, imageType ImageType) (*Image, error) {
	open, err := os.Open(fPath)
	if err != nil {
		return nil, err
	}
	//defer open.Close()

	dirs := strings.Split(filepath.Dir(fPath), "/")
	blogName := dirs[len(dirs)-1]
	srcImg := &Image{
		content:     open,
		imageFormat: imageType,
		blog:        blogName,
	}
	return srcImg, nil
}

type ImageHandler func(image Image)

//CopyImgFromDirOfBlogToDirImgCDN copy images of blogs directory to image and process it if `imageHandler` function is not empty;
//https://opensource.com/article/18/6/copying-files-go
func (o *SyncImageOption) CopyImgFromDirOfBlogToDirImgCDN(srcDir, targetDir string, imageHandlers ...ImageHandler) error {
	logrus.Info("🚀: sync images from blog directory to github cdn directory")
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		ok, imageType := IsImage(path)
		if !info.IsDir() && ok {
			srcImg, err := NewImage(path, imageType)
			if err != nil {
				logrus.Warn(err)
				return nil
			}
			defer srcImg.content.Close()

			// Image will be compressed before copied to target dir
			if o.CompressEnable {
				if info.Size() > MaxImageSize {
					// compress image
					logrus.Infof("♻️: compress image:%s", srcImg.NameWithBlog())
					err := CompressImageByCommandTool(*srcImg)
					if err != nil {
						fmt.Println()
						return err
					}
					// FIXME:whether reopen the image which has been compressed.
				}
			}

			// check if blog directory exists
			blogDir := filepath.Join(targetDir, srcImg.blog)
			if !dirExists(blogDir) {
				err := os.MkdirAll(blogDir, 0777)
				if err != nil {
					logrus.Errorf("failed create blog dir :%s", err)
					return nil
				}
			}

			targetImage := filepath.Join(targetDir, srcImg.NameWithBlog())
			if !fileExists(targetImage) || o.Update {
				logrus.Infof("🔨: copy image :%s",srcImg.NameWithBlog())
				_ = os.Remove(targetImage)

				targetImgFile, err := os.Create(targetImage)
				if err != nil {
					logrus.Warn(err)
					return nil
				}
				defer targetImgFile.Close()
				_, err = io.Copy(targetImgFile, srcImg.content)
				if err != nil {
					logrus.Warn(err)
					return nil
				}

				if srcImg.imageFormat != WEBP {
					logrus.Infof("🔋: generate to webp :%s",targetImage)
					GenerateWebpFormat(targetImage)
				}
			}

			if len(imageHandlers) != 0 {
				for _, handler := range imageHandlers {
					handler(*srcImg)
				}
			}
			//if err = os.Chdir(filepath.Dir(targetImage));err!= nil {
			//
			//}
			logrus.Infoln()
		}
		return nil
	})
	return err
}

//ReleaseImageTOCDN
func ReleaseImageTOCDN(imageCDNRootDir string) {
	command := exec.Command("/bin/bash", "-c", `git add . && git commit -m "sync" && release-it --ci`)
	command.Dir = imageCDNRootDir
	command.Stdout = os.Stdout
	command.Run()
}

//GenerateWebpFormat generate a webp format image of origin image and the generated image will be placed in same dir.
func GenerateWebpFormat(imagePath string) {
	command := exec.Command("cwebp", imagePath, "-o", strings.TrimSuffix(imagePath, filepath.Ext(imagePath))+".webp")
	logrus.Debugf("the directory where command executed :%s,the command string :%s", command.Dir, command.String())
	command.Stdout = os.Stdout
	command.Run()
}

func PrintErrorToStdErr(err error, message string) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "%s: %s", message, err)
	os.Exit(1)
}

//CompressImageByTinyPNGAPI compress image by the tinypng api.
func CompressImageByTinyPNGAPI(image Image) {

}

//CompressImageByCommandTool compress image by the local command tool.
func CompressImageByCommandTool(image Image) error {
	var tool string
	switch image.imageFormat {
	case PNG:
		tool = "optipng"
	case JPEG:
		tool = "jpegtopnm"
	default:
		tool = ""
	}
	if tool == "" {
		logrus.Errorf("not support image type:%s", image.imageFormat)
		return errors.New("nonsupport image type")
	}
	cmd := exec.Command(tool, image.Path())
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

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

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// dirExists checks if a directory exists and is a directory before we
// try using it to prevent further errors.
func dirExists(dirname string) bool {
	info, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
