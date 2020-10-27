package blog

import (
	"fmt"
	"hack/cmd/util"
	"hack/pkg"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

type PublishOption struct {
	SyncToBlogOriginRepo bool
}

func NewPublishOption() *PublishOption {
	return &PublishOption{}
}

func NewPublishCommand(cfg util.BlogConfig) *cobra.Command {
	o := NewPublishOption()

	cmd := &cobra.Command{
		Use:   "publish",
		Short: "publish blog to internet",
		Long:  "publish blog to internet",
		Run: func(cmd *cobra.Command, args []string) {
			o.Run(cfg.BlogDir)
		},
	}
	cmd.Flags().BoolVar(&o.SyncToBlogOriginRepo, "sync", true, "if flag `sync` is set true,the blog directory change will be pushed to default remote repo")
	return cmd
}

func (o *PublishOption) Run(blogDir string) {
	if o.SyncToBlogOriginRepo {
		promptString, err := pkg.PromptString("Message")
		if err != nil {
			fmt.Println(err)
		}
		err = o.SyncToGithubRepo(blogDir, promptString)
		PrintErrorToStdErr(err,"error happened")
	}
	o.Deploy(blogDir)
}

func (o *PublishOption) Deploy(blogDir string) {
	command := exec.Command("/bin/bash", "-c", "hexo generate -f && gulp && hexo deploy")
	command.Dir = blogDir
	command.Stdout = os.Stdout
	command.Run()
}

func (o *PublishOption) SyncToGithubRepo(blogDir, message string) error{
	command := exec.Command("/bin/bash", "-c", fmt.Sprintf(`git add . && git commit -m "%s" && git push`, message))
	command.Dir = blogDir
	command.Stdout = os.Stdout
	return command.Run()
}
