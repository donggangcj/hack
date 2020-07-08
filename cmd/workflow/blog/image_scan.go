package blog

import (
	"hack/cmd/util"
	"hack/pkg"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewImageScan(cfg util.BlogConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "list unsynchronized image",
		Long:  "list unsynchronized image",
		Run: func(cmd *cobra.Command, args []string) {
			RunImageScan(cfg)
		},
	}
	return cmd
}

func RunImageScan(cfg util.BlogConfig) {
	images := ListDraftImage(cfg.BlogSourceDir, cfg.ImageRepoDir)
	i := pkg.Images(images)
	i.Print(os.Stdout)
}

//ListDraftImage list unsynchronized images
func ListDraftImage(srcDir, targetDir string) []pkg.Image {
	images := make([]pkg.Image, 0)
	filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
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
			if !fileExists(targetImage) {
				images = append(images, *srcImg)
			}
		}
		return nil
	})
	return images
}
