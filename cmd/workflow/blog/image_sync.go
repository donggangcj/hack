package blog

import (
	"context"
	"fmt"
	"hack/cmd/util"
	"hack/pkg"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//MaxImageSize present the max size that a image need to be compress
const MaxImageSize = 1024000

type SyncImageOption struct {
	CompressEnable  bool
	ReleaseEnable   bool
	TinyPNGEnable   bool
	AsciinemaEnable bool
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
	cmd.Flags().BoolVar(&o.TinyPNGEnable, "tinypng", false, "compress image by tinypng,the image compressed by this way has smaller size but it needs more time")
	cmd.Flags().BoolVar(&o.AsciinemaEnable, "asciinema", false, "sync asciinema files of blog")
	return cmd
}

func (o *SyncImageOption) Run(cfg util.BlogConfig) {
	err := o.CopyImgFromDirOfBlogToDirImgCDN(cfg.BlogSourceDir, cfg.ImageRepoDir, cfg.TinyPNGToken)
	PrintErrorToStdErr(err, "err happen when sync image")

	if o.ReleaseEnable {
		logrus.Info("🚀: release image")
		ReleaseImageTOCDN(cfg.ImageRepoRootDir)
	}
}

type ImageHandler func(image pkg.Image)

//CopyImgFromDirOfBlogToDirImgCDN copy images of blogs directory to image and process it if `imageHandler` function is not empty;
//https://opensource.com/article/18/6/copying-files-go
func (o *SyncImageOption) CopyImgFromDirOfBlogToDirImgCDN(srcDir, targetDir string, token string, imageHandlers ...ImageHandler) error {
	logrus.Info("🚀: sync images from blog directory to github cdn directory")
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		ok, imageType := pkg.IsImage(path)
		if !info.IsDir() && ok {
			srcImg, err := pkg.NewImage(path, imageType)
			if err != nil {
				logrus.Warn(err)
				return nil
			}

			// check if blog directory exists
			blogDir := filepath.Join(targetDir, srcImg.Blog)
			if !dirExists(blogDir) {
				err := os.MkdirAll(blogDir, 0777)
				if err != nil {
					logrus.Errorf("failed create blog dir :%s", err)
					return nil
				}
			}

			targetImage := filepath.Join(targetDir, srcImg.NameWithBlog())
			if !fileExists(targetImage) || o.Update {
				// Image will be compressed before copied to target dir
				if o.CompressEnable {
					if info.Size() > MaxImageSize {
						// compress image
						logrus.Infof("♻️: compress image:%s", srcImg.NameWithBlog())
						var err error
						if o.TinyPNGEnable {
							err = util.CompressImageByTinyPNGAPI(context.Background(), srcImg.Path, token)
						} else {
							err = pkg.CompressImageByCommandTool(*srcImg)
						}
						if err != nil {
							fmt.Println()
							return err
						}
						// FIXME:whether reopen the image which has been compressed.
					}
				}

				logrus.Infof("🔨: copy image :%s", srcImg.NameWithBlog())
				err := CopyFile(srcImg.Path, targetImage)
				if err != nil {
					logrus.Errorf("failed copy file from %s to %s: %s", srcImg.Path, targetImage, err)
					return nil
				}

				if srcImg.ImageFormat != pkg.WEBP {
					logrus.Infof("🔋: generate to webp :%s", targetImage)
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
		}
		if o.AsciinemaEnable {
			if IsAsciinema(path) {
				dirs := strings.Split(filepath.Dir(path), "/")
				blogName := dirs[len(dirs)-1]

				blogDir := filepath.Join(targetDir, blogName)
				if !dirExists(blogDir) {
					err := os.MkdirAll(blogDir, 0777)
					if err != nil {
						logrus.Errorf("failed create blog dir :%s", err)
						return nil
					}
				}

				targetFile := filepath.Join(blogDir, filepath.Base(path))
				if !fileExists(targetFile) {
					logrus.Infof("🔨: copy asciinema :%s", filepath.Join(blogName, filepath.Base(path)))
					err := CopyFile(path, targetFile)
					if err != nil {
						logrus.Errorf("failed copy asciinema file: %s", err)
						return nil
					}
				}

			}
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

func IsAsciinema(fPath string) bool {
	return strings.HasSuffix(fPath, ".cast")
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

func CopyFile(sourceFile, targetFile string) error {
	_ = os.Remove(targetFile)

	open, err := os.Open(sourceFile)
	if err != nil {
		return nil
	}
	defer open.Close()

	targetImgFile, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer targetImgFile.Close()
	_, err = io.Copy(targetImgFile, open)
	if err != nil {
		return err
	}
	return nil
}
