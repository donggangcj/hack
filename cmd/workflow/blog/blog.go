package blog

import (
	"hack/cmd/util"

	"github.com/spf13/cobra"
)

func NewBlogCommand(cfg util.BlogConfig) *cobra.Command {
	validArgs := []string{"asciinema", "image", "publish", "sync"}
	cmd := &cobra.Command{
		Use:   "blog",
		Short: "toolsets which help manage the blog",
		Long:  "toolsets which help manage the blog",
		ValidArgs: validArgs,
	}
	cmd.AddCommand(NewSyncImage(cfg))
	cmd.AddCommand(NewPublishCommand(cfg))
	cmd.AddCommand(NewImageScan(cfg))
	cmd.AddCommand(NewUnsynchronizedAsciinemas(cfg))
	return cmd
}
