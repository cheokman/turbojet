package heater

import (
	"fmt"
	"time"
	"turbojet/cli"
)

type Heater struct {
	cNode map[string][]string
	files []string
	path  string
	typ   string
}

func NewHeater(cNode map[string][]string, path string, files []string, typ string) *Heater {
	return &Heater{
		cNode: cNode,
		files: files,
		path:  path,
		typ:   typ,
	}
}

func (h *Heater) Process(ctx *cli.Context) []*Job {
	var jobs []*Job
	w := ctx.Writer()

	for domain, ips := range h.cNode {
		j := NewJob(ips, domain, h.path, h.files)
		cli.Printf(w, "Start Job for domain[%s]\n", domain)
		start := time.Now()
		j.Process(ctx)
		processTime := time.Since(start)
		js := j.GetSummary()
		for _, ips := range js.IPSummaries {
			fmt.Printf("%s Processing Time: %d, Download: %d\n", ips.IP, ips.TotalProcessTime, ips.TotalDownload)
		}
		cli.Printf(w, "Done Job for domain[%s] in %s\n", domain, processTime)
		jobs = append(jobs, j)
	}
	return jobs
}
