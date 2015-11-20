package filesystems

import (
	"fmt"
	"github.com/Symantec/Dominator/lib/format"
	"io"
	"sort"
)

type fileSystemsList []*fileSystem

func (p *prober) writeHtml(writer io.Writer) {
	fsList := make(fileSystemsList, 0, len(p.fileSystems))
	for _, fs := range p.fileSystems {
		fsList = append(fsList, fs)
	}
	sort.Sort(fsList)
	fmt.Fprintln(writer, `<style>
                          table, th, td {
                          border-collapse: collapse;
                          }
                          </style>`)
	//fmt.Fprintln(writer, `<table border="1" style="width:100%">`)
	fmt.Fprintln(writer, `<table border="1">`)
	fmt.Fprintln(writer, "  <tr>")
	fmt.Fprintln(writer, "    <th>Mount Point</th>")
	fmt.Fprintln(writer, "    <th>Size</th>")
	fmt.Fprintln(writer, "    <th>Used</th>")
	fmt.Fprintln(writer, "    <th></th>")
	fmt.Fprintln(writer, "  </tr>")
	for _, fs := range fsList {
		fs.writeHtml(writer)
	}
	fmt.Fprintln(writer, "</table>")
	fmt.Fprintln(writer, "</body>")
}

func (fsList fileSystemsList) Len() int {
	return len(fsList)
}

func (fsList fileSystemsList) Less(left, right int) bool {
	return fsList[left].mountPoint < fsList[right].mountPoint
}

func (fsList fileSystemsList) Swap(left, right int) {
	fsList[left], fsList[right] = fsList[right], fsList[left]
}

func (fs *fileSystem) writeHtml(writer io.Writer) {
	usedBytes := fs.size - fs.free
	usedPercent := float64(usedBytes) * 100 / float64(fs.size)
	fmt.Fprintf(writer, "  <tr>\n")
	fmt.Fprintf(writer, "    <td><center>%s</td>\n", fs.mountPoint)
	fmt.Fprintf(writer, "    <td><center>%s</td>\n", format.FormatBytes(fs.size))
	fmt.Fprintf(writer, "    <td><center>%.1f%%</td>\n", usedPercent)
	fmt.Fprint(writer, "    <td>")
	fs.writeHtmlBar(writer)
	fmt.Fprintln(writer, "</td>")
	fmt.Fprintln(writer, "  </tr>")
}

func (fs *fileSystem) writeHtmlBar(writer io.Writer) {
	//usedBytes := fs.size - fs.free
	bgColour := "green"
	if fs.available < 1 {
		bgColour = "red"
	}
	fmt.Fprint(writer, `<table border="0" style="width:200px">`)
	fmt.Fprintf(writer,
		"<tr><td style=\"width:100%%;background-color:%s\">&nbsp;</td></tr>",
		bgColour)
	fmt.Fprint(writer, "</table>")
}
