package workflow

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

//ProxyOption
type ProxyOption struct {
	On bool
}

func NewCmdProxy() *cobra.Command {
	o := &ProxyOption{}
	cmd := &cobra.Command{
		Use:   "proxy ",
		Short: "open or close ShadowsocksX-NG-R8 and Proxifier at same time,make sure those applications have been installed",
		Long:  "open or close ShadowsocksX-NG-R8 and Proxifier at same time",
		Run: func(cmd *cobra.Command, args []string) {
			if o.On {
				OpenProxySets()
				return
			}else{
				CloseProxySet()
			}
		},
	}
	cmd.Flags().BoolVar(&o.On, "on", true, "")
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
}

func CloseProxySet() {
	// close proxifier
	command := exec.Command("osascript","-e",`quit app "Proxifier"`)
	output, err := command.Output()
	if err != nil {
		fmt.Fprintf(os.Stdout,"error:%s",output)
		return
	}
	// close ShadowsocksX-NG-R8
	command = exec.Command("osascript","-e",`quit app "ShadowsocksX-NG-R8"`)
	output, err = command.Output()
	if err != nil {
		fmt.Fprintf(os.Stdout,"error:%s",output)
		return
	}
}

