package workflow

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCmdProxy() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proxy ",
		Short: "open or close ShadowsocksX-NG-R8 and Proxifier at same time",
		Long:  "open or close ShadowsocksX-NG-R8 and Proxifier at same time",
	}
	cmd.AddCommand(newProxyOpen())
	cmd.AddCommand(newProxyClose())
	return cmd
}

func newProxyOpen() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open",
		Short: "open ShadowsocksX-NG-R8 and Proxifier",
		Long:  "open ShadowsocksX-NG-R8 and Proxifier",
		Run: func(cmd *cobra.Command, args []string) {
			OpenProxySets()
		},
	}
	return cmd
}

func OpenProxySets() {
	// open shadowsocksX-NG-R8
	command := exec.Command("open", "-a", `ShadowsocksX-NG-R8`)
	output, err := command.Output()
	if err != nil {
		fmt.Fprintf(os.Stdout, "error:%s", output)
		return
	}
	// open proxifier
	command = exec.Command("open", "-a", "Proxifier")
	output, err = command.Output()
	if err != nil {
		fmt.Fprintf(os.Stdout, "error:%s", output)
		return
	}
	logrus.Info("🔈 proxy is on")
}

func newProxyClose() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close",
		Short: "close ShadowsocksX-NG-R8 and Proxifier",
		Long:  "close ShadowsocksX-NG-R8 and Proxifier",
		Run: func(cmd *cobra.Command, args []string) {
			CloseProxySet()
		},
	}
	return cmd
}
func CloseProxySet() {
	// close proxifier
	command := exec.Command("osascript", "-e", `quit app "Proxifier"`)
	output, err := command.Output()
	if err != nil {
		fmt.Fprintf(os.Stdout, "error:%s", output)
		return
	}
	// close ShadowsocksX-NG-R8
	command = exec.Command("osascript", "-e", `quit app "ShadowsocksX-NG-R8"`)
	output, err = command.Output()
	if err != nil {
		fmt.Fprintf(os.Stdout, "error:%s", output)
		return
	}
	logrus.Info("🔈 proxy is off")
}
