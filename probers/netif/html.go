package netif

import (
	"fmt"
	"io"
	"sort"
)

type netInterfacesList []*netInterface

func (p *prober) writeHtml(writer io.Writer) {
	names := make([]string, 0, len(p.netInterfaces))
	for name, netIf := range p.netInterfaces {
		if netIf.speed > 0 {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	netIfList := make(netInterfacesList, 0, len(p.netInterfaces))
	for _, name := range names {
		netIfList = append(netIfList, p.netInterfaces[name])
	}
	fmt.Fprintln(writer, `<style>
                          table, th, td {
                          border-collapse: collapse;
                          }
                          </style>`)
	fmt.Fprintln(writer, `<table border="1">`)
	fmt.Fprintln(writer, "  <tr>")
	fmt.Fprintln(writer, "    <th>Name</th>")
	fmt.Fprintln(writer, "    <th>Rx</th>")
	fmt.Fprintln(writer, "    <th>Utilisation</th>")
	fmt.Fprintln(writer, "    <th>Tx</th>")
	fmt.Fprintln(writer, "    <th>Utilisation</th>")
	fmt.Fprintln(writer, "  </tr>")
	for _, fs := range netIfList {
		fs.writeHtml(writer)
	}
	fmt.Fprintln(writer, "</table>")
	fmt.Fprintln(writer, "</body>")
}

func (netIf *netInterface) writeHtml(writer io.Writer) {
	rxUtilisation := float64(netIf.rxDataRate) / float64(netIf.speed)
	txUtilisation := float64(netIf.txDataRate) / float64(netIf.speed)
	fmt.Fprintln(writer, "  <tr>")
	fmt.Fprintf(writer, "    <td><center>%s</td>\n", netIf.name)
	fmt.Fprintf(writer, "    <td><center>%.1f%%</td>\n", rxUtilisation*100)
	writeHtmlBar(writer, rxUtilisation)
	fmt.Fprintf(writer, "    <td><center>%.1f%%</td>\n", txUtilisation*100)
	writeHtmlBar(writer, txUtilisation)
	fmt.Fprintln(writer, "  </tr>")
}

func writeHtmlBar(writer io.Writer, utilisation float64) {
	barColour := "green"
	if utilisation >= 0.95 {
		barColour = "red"
	} else if utilisation >= 0.9 {
		barColour = "orange"
	} else if utilisation >= 0.75 {
		barColour = "yellow"
	}
	fmt.Fprint(writer, `<td><table border="0" style="width:100px"><tr>`)
	fmt.Fprintf(writer,
		"<td style=\"width:%.1f%%;background-color:%s\">&nbsp;</td>",
		utilisation*100, barColour)
	fmt.Fprintf(writer,
		"<td style=\"width:%.1f%%\"></td>", 100-utilisation*100)
	fmt.Fprint(writer, "</tr></table></td>")
}
