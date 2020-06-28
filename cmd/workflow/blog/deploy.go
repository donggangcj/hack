package blog

import (
	"hack/cmd/util"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

type PublishOption struct {
}

func NewPublishOption() *PublishOption {
	return &PublishOption{}
}

func NewPublishCommand(cfg util.BlogConfig) *cobra.Command {
	o := NewPublishOption()

	cmd := &cobra.Command{
		Use:   "publish",
		Short: "publish image to internet",
		Long:  "publish image to internet",
		Run: func(cmd *cobra.Command, args []string) {
			o.Run(cfg.BlogDir)
		},
	}
	return cmd
}

func (o *PublishOption) Run(blogDir string) {
	o.Deploy(blogDir)
}

func (o *PublishOption) Deploy(blogDir string) {
	command := exec.Command("/bin/bash", "-c", "hexo generate -f && gulp && hexo deploy")
	command.Dir = blogDir
	command.Stdout = os.Stdout
	command.Run()
}

func (o *PublishOption) SyncToGithubRepo(blogDir string) {

}

