package workflow

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ProxyOption struct {
	ssr bool
}

func NewCmdProxy() *cobra.Command {
	o := ProxyOption{}
	cmd := &cobra.Command{
		Use:   "proxy ",
		Short: "Open proxy software",
		Long:  "Close proxy software",
	}
	cmd.PersistentFlags().BoolVar(&o.ssr, "ssr", false, "use SSR to replace V2ray")
	cmd.AddCommand(newProxyOpen(o))
	cmd.AddCommand(newProxyClose(o))
	return cmd
}

func newProxyOpen(option ProxyOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open",
		Short: "open ShadowsocksX-NG-R8 and Proxifier",
		Long:  "open ShadowsocksX-NG-R8 and Proxifier",
		Run: func(cmd *cobra.Command, args []string) {
			OpenProxySets(option.ssr)
		},
	}
	return cmd
}

func OpenProxySets(ssr bool) {
	// open shadowsocksX-NG-R8
	proxyType := "V2rayU"
	if ssr  {
		proxyType = "ShadowsocksX-NG-R8"
	}
	command := exec.Command("open", "-a", proxyType)
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

func newProxyClose(option ProxyOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close",
		Short: "close ShadowsocksX-NG-R8 and Proxifier",
		Long:  "close ShadowsocksX-NG-R8 and Proxifier",
		Run: func(cmd *cobra.Command, args []string) {
			CloseProxySet(option.ssr)
		},
	}
	return cmd
}
func CloseProxySet(ssr bool) {
	// close proxifier
	command := exec.Command("osascript", "-e", `quit app "Proxifier"`)
	output, err := command.Output()
	if err != nil {
		fmt.Fprintf(os.Stdout, "error:%s", output)
		return
	}
	// close ShadowsocksX-NG-R8
	proxyType := "V2rayU"
	if ssr  {
		proxyType = "ShadowsocksX-NG-R8"
	}
	command = exec.Command("osascript", "-e", fmt.Sprintf(`quit app "%s"`, proxyType))
	output, err = command.Output()
	if err != nil {
		fmt.Fprintf(os.Stdout, "error:%s", output)
		return
	}
	logrus.Info("🔈 proxy is off")
}
