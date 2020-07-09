package blog

import (
	"hack/cmd/util"
	"hack/pkg"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewUnsynchronizedAsciinemas(cfg util.BlogConfig) *cobra.Command {
	cobra := &cobra.Command{
		Use:   "asciinema",
		Short: "list unsynchronized asciinema",
		Long:  "list unsynchronized asciinema",
		Run: func(cmd *cobra.Command, args []string) {
			RunUnsynchronized(cfg)
		},
	}
	return cobra
}

func RunUnsynchronized(cfg util.BlogConfig) {
	asciinemas := ListDraftAsciinemas(cfg.BlogSourceDir, cfg.ImageRepoDir)
	a := pkg.Asciinemas(asciinemas)
	a.Print(os.Stdout)
}

//ListDraftImage list unsynchronized images
func ListDraftAsciinemas(srcDir, targetDir string) []pkg.Asciinema {
	images := make([]pkg.Asciinema, 0)
	filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		ok :=pkg.IsAsciinema(path)
		if !info.IsDir() && ok {
			srcAscii, err := pkg.NewAsciinema(path)
			if err != nil {
				logrus.Warn(err)
				return nil
			}

			// check if blog directory exists
			blogDir := filepath.Join(targetDir, srcAscii.Blog)
			if !dirExists(blogDir) {
				err := os.MkdirAll(blogDir, 0777)
				if err != nil {
					logrus.Errorf("failed create blog dir :%s", err)
					return nil
				}
			}
			targetAsciinema := filepath.Join(blogDir, filepath.Base(srcAscii.Path))
			if !fileExists(targetAsciinema) {
				images = append(images, *srcAscii)
			}
		}
		return nil
	})
	return images
}

