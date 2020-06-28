package blog

import (
	"fmt"
	"hack/cmd/util"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func NewBlogCommand(cfg util.BlogConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blog",
		Short: "toolsets which help manage the blog",
		Long:  "toolsets which help manage the blog",
		Run: func(cmd *cobra.Command, args []string) {
			command := exec.Command("hexo-image-sync", "-v", "--release")
			command.Stdout = os.Stdout
			err := command.Run()
			if err != nil {
				fmt.Fprintf(os.Stdout, "error:%s", err)
			}
		},
	}
	cmd.AddCommand(NewSyncImage(cfg))
	cmd.AddCommand(NewPublishCommand(cfg))
	cmd.AddCommand(NewImageScan(cfg))
	return cmd
}
