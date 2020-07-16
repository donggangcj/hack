package cmd

import (
	"fmt"
	"hack/cmd/util"
	"hack/cmd/workflow"
	"hack/cmd/workflow/blog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hack",
	Short: "hack is a toolset to simplify workflow ",
	Long:  "hack is a toolset to simplify workflow, hack reads config from `$HOME/.hackconfig`, so you should make sure that file existed",
}

func Register() {
	dir, _ := os.UserHomeDir()
	cfg, err := util.BuildConfigFromFile(filepath.Join(dir,".hackconfig.yaml"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read config from configfile:%s", err.Error())
		os.Exit(1)
	}

	rootCmd.AddCommand(workflow.NewCmdProxy())
	rootCmd.AddCommand(blog.NewBlogCommand(cfg.BlogConfig))
	rootCmd.AddCommand(NewCompletionCommand())
}

func Execute() {
	Register()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
