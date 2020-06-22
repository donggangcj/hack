package blog

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

type PublishOption struct {
	blogDir string
}

func NewPublishOption() *PublishOption {
	return &PublishOption{}
}

func NewPublishCommand() *cobra.Command {
	o := NewPublishOption()

	cmd := &cobra.Command{
		Use:   "publish",
		Short: "publish image to internet",
		Long:  "publish image to internet",
		Run: func(cmd *cobra.Command, args []string) {
			o.Run()
		},
	}
	return cmd
}

func (o *PublishOption) Run() {
	command := exec.Command("/bin/bash", "-c", "hexo clean & hexo deploy")
	command.Dir = o.blogDir
	command.Stdout = os.Stdout
	command.Run()
}
