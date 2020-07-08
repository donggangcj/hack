package pkg

import (
	"hack/cmd/util"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type Asciinema struct {
	Path string
	Size int64
	Blog string
}

func (a *Asciinema) ToTableColumnString() []string{
	return []string{a.Blog, filepath.Base(a.Path), util.ByteCountSI(a.Size)}
}

func NewAsciinema(fPath string) (*Asciinema, error) {
	stat, err := os.Stat(fPath)
	if err != nil {
		return nil, err
	}
	dirs := strings.Split(filepath.Dir(fPath), "/")
	blogName := dirs[len(dirs)-1]
	return &Asciinema{
		Path: fPath,
		Size: stat.Size(),
		Blog: blogName,
	}, nil
}

type Asciinemas []Asciinema

func (a Asciinemas) Print(writer io.Writer) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Blog", "Asciinema", "Size"})
	table.SetRowLine(true)
	table.SetAutoMergeCells(true)
	//table.SetColumnColor(tablewriter.Colors{tablewriter.BgGreenColor}, tablewriter.Colors{tablewriter.BgBlueColor},
	//	tablewriter.Colors{tablewriter.BgHiRedColor})
	for _, image := range a {
		table.Append(image.ToTableColumnString())
	}
	//table.SetFooter([]string{"", "Total", strconv.Itoa(len(i))})
	table.Render()
}

