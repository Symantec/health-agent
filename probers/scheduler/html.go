package scheduler

import (
	"fmt"
	"io"
)

func (p *prober) writeHtml(writer io.Writer) {
	fmt.Fprintln(writer, `<style>
                          table, th, td {
                          border-collapse: collapse;
                          }
                          </style>`)
	fmt.Fprint(writer, "Load average: ")
	p.writeLoadAverage(writer, "1 min", p.loadavg.oneMinute)
	fmt.Fprint(writer, " / ")
	p.writeLoadAverage(writer, "5 min", p.loadavg.fiveMinutes)
	fmt.Fprint(writer, " / ")
	p.writeLoadAverage(writer, "15 min", p.loadavg.fifteenMinutes)
	fmt.Fprintln(writer, "<br>")
	fmt.Fprintln(writer, `<table border="1">`)
	fmt.Fprintln(writer, "  <tr>")
	fmt.Fprintln(writer, "    <th>CPU Time</th>")
	fmt.Fprintln(writer, "    <th>Value</th>")
	fmt.Fprintln(writer, "    <th></th>")
	fmt.Fprintln(writer, "  </tr>")
	writeMetric(writer, "User", p.cpuStats.userTime.fraction)
	writeMetric(writer, "User Niced", p.cpuStats.userNiceTime.fraction)
	writeMetric(writer, "System", p.cpuStats.systemTime.fraction)
	fmt.Fprintln(writer, "</table>")
	fmt.Fprintln(writer, "</body>")
}

func (p *prober) writeLoadAverage(writer io.Writer, period string,
	value float32) {
	loadavgColour := "green"
	if value > float32(p.numCpus*2) {
		loadavgColour = "red"
	} else if value > float32(p.numCpus) {
		loadavgColour = "orange"
	}
	fmt.Fprintf(writer, "<font color=\"%s\">%.2f</font> (%s)",
		loadavgColour, value, period)
}

func writeMetric(writer io.Writer, name string, utilisation float64) {
	barColour := "green"
	if utilisation >= 0.95 {
		barColour = "red"
	} else if utilisation >= 0.9 {
		barColour = "orange"
	} else if utilisation >= 0.75 {
		barColour = "yellow"
	}
	fmt.Fprintln(writer, "  <tr>")
	fmt.Fprintf(writer, "    <td><center>%s</td>\n", name)
	fmt.Fprintf(writer, "    <td><center>%.1f%%</td>\n", utilisation*100)
	fmt.Fprint(writer, `    <td><table border="0" style="width:100px"><tr>`)
	fmt.Fprintf(writer,
		"<td style=\"width:%.1f%%;background-color:%s\">&nbsp;</td>",
		utilisation*100, barColour)
	fmt.Fprintf(writer,
		"<td style=\"width:%.1f%%\"></td>", 100-utilisation*100)
	fmt.Fprintln(writer, "</tr></table></td>")
	fmt.Fprintln(writer, "  </tr>")
}
