package proberlist

import (
	"fmt"
	"github.com/Symantec/health-agent/lib/prober"
	"io"
	"net/http"
	"sort"
	"strings"
)

type dashboardType struct {
	name             string
	dashboardYielder prober.DashboardYielder
}

func (pl *ProberList) writeHtml(writer io.Writer, req *http.Request) {
	dashboards := make([]dashboardType, 0)
	for _, p := range pl.getProbers() {
		if htmler, ok := p.prober.(HtmlWriter); ok {
			htmler.WriteHtml(writer)
			fmt.Fprintln(writer, "<br>")
		}
		if dashboardYielder, ok := p.prober.(prober.DashboardYielder); ok {
			dashboards = append(dashboards,
				dashboardType{p.name, dashboardYielder})
		}
	}
	if len(dashboards) < 1 {
		return
	}
	sort.Slice(dashboards, func(i, j int) bool {
		return dashboards[i].name < dashboards[j].name
	})
	protocol := "http"
	if req.TLS != nil {
		protocol = "https"
	}
	host := strings.Split(req.Host, ":")[0]
	for _, dashboard := range dashboards {
		colour := "green"
		if !dashboard.dashboardYielder.HealthCheck() {
			colour = "red"
		}
		port := dashboard.dashboardYielder.GetPort()
		fmt.Fprintf(writer,
			"<font color=\"%s\">%s</font> <a href=\"%s://%s:%d\">dashboard</a><br>\n",
			colour, dashboard.name, protocol, host, port)
	}
	fmt.Fprintln(writer, "<br>")
}
