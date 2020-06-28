package blog

import (
	"hack/cmd/util"
	"os"

	"github.com/spf13/cobra"
)

func NewImageScan(cfg util.BlogConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "list unsynchronized image",
		Long:  "list unsynchronized image",
		Run: func(cmd *cobra.Command, args []string) {
			Run(cfg)
		},
	}
	return cmd
}

func Run(cfg util.BlogConfig) {
	images := ListDraftImage(cfg.BlogSourceDir, cfg.ImageRepoDir)
	i := Images(images)
	i.Print(os.Stdout)
}
