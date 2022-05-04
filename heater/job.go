package heater

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
	"turbojet/cli"

	"github.com/imdario/go-ulid"
)

type Job struct {
	ID       string
	domain   string
	hostname string
	port     string
	schema   string
	path     string
	files    []string
	ips      []string
	summary  *JobSummary
}

type JobSummary struct {
	*Summary
	IPSummaries []*IPSummary
}

func NewJob(ips []string, domain string, path string, files []string) *Job {
	u, err := url.Parse(domain)
	if err != nil {
		return &Job{}
	}
	hostname, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		hostname = u.Host
		port = "80"
	}
	return &Job{
		ID:       ulid.New().String(),
		domain:   domain,
		schema:   u.Scheme,
		hostname: hostname,
		port:     port,
		path:     path,
		files:    files,
		ips:      ips,
		summary:  &JobSummary{},
	}
}

func (j *Job) Process(ctx *cli.Context) (*JobSummary, error) {
	w := ctx.Writer()
	done := make(chan interface{})
	taskChan := make(chan *Task)
	defer close(done)

	tasks, err := j.createTasks(j.ips, j.hostname, j.path, j.files)
	if err != nil {
		cli.Printf(w, "Job create tasks error: %s\n", err)
		return &JobSummary{}, err
	}

	cli.Printf(w, "Total %d of heat tasks created\n", len(tasks))
	cli.Printf(w, "starting job task process\n")
	taskSummary := j.processTask(done, taskChan)

	cli.Printf(w, "started job task process\n")
	go func() {
		defer close(taskChan)
		totalTask := len(tasks)
		startTime := time.Now()
		for i, t := range tasks {
			if i > 0 && i%500 == 0 {
				processTime := time.Since(startTime)
				cli.Printf(w, "send %d/%d task in %s\n", i, totalTask, processTime)
				startTime = time.Now()
			}
			taskChan <- t
		}
		cli.Print(w, "All tasks are dispatched and waiting to complete\n")
	}()

	for ts := range taskSummary {
		j.PutTaskSummary(ts)
	}
	return j.GetSummary(), nil
}

func (j *Job) heat(t *Task) error {
	domain := t.GetHostname()
	url := t.GetURL()
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Timeout: 60 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.PutError(&err)
		fmt.Println("Error:", err)
		return err
	}

	req.Host = domain

	resp, err := client.Do(req)
	if err != nil {
		t.PutError(&err)
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.PutError(&err)
		return err
	}
	t.PutDownloadSize(len(body))
	t.PutResponse(resp)
	return nil
}

func (j *Job) processTask(done <-chan interface{},
	taskChan <-chan *Task,
) <-chan TaskSummary {
	var wg sync.WaitGroup
	respChan := make(chan TaskSummary)

	processor := func(c <-chan *Task) {
		defer wg.Done()
		for t := range c {
			// fmt.Printf("Get URL:%s \n", t.GetURL())
			start := time.Now()
			err := j.heat(t)
			if err != nil {
				fmt.Printf("Processing heat task error: %s\n", err)
			}
			pt := time.Since(start)
			t.PutProcessTime(pt)
			// fmt.Printf("Done URL:%s in %s processing time\n", t.GetURL(), pt)
			select {
			case <-done:
				return
			case respChan <- t.GetSummary():
			}
		}
	}

	numWorkers := 100 // runtime.NumCPU()
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go processor(taskChan)
	}

	go func() {
		wg.Wait()
		close(respChan)
	}()
	return respChan
}

func (j *Job) createTasks(ips []string, domain string, path string, files []string) ([]*Task, error) {
	var tasks []*Task
	for _, ip := range ips {
		for _, f := range files {
			t := NewTask(ip, j.schema, j.hostname, path, f)
			tasks = append(tasks, t)
		}
	}
	return tasks, nil
}

func (j *Job) PutTaskSummary(ts TaskSummary) error {
	j.putSummaryByIP(ts.IP, ts)
	return nil
}

func (j *Job) GetSummary() *JobSummary {
	for _, ips := range j.summary.IPSummaries {
		ips.GetStat()
	}
	return j.summary
}

func (j *Job) putSummaryByIP(ip string, ts TaskSummary) {
	js := j.summary
	for _, ips := range js.IPSummaries {
		if ip == ips.IP {
			ips.PutSummary(ts)
			return
		}
	}
	newIPS := &IPSummary{
		IP: ip,
	}
	newFileS := NewFileSummary(ts.File, ts.DownloadSize, ts.ProcessTime, ts.Err)
	newIPS.FileSummaries = append(newIPS.FileSummaries, newFileS)
	js.IPSummaries = append(js.IPSummaries, newIPS)
}

type IPSummary struct {
	IP string
	Summary
	FileSummaries []*FileSummary
}

func (ips *IPSummary) PutSummary(ts TaskSummary) {
	newFS := NewFileSummary(ts.File, ts.DownloadSize, ts.ProcessTime, ts.Err)
	for i, fs := range ips.FileSummaries {
		if fs.File == ts.File {
			ips.FileSummaries[i] = newFS
			return
		}
	}
	ips.FileSummaries = append(ips.FileSummaries, newFS)
}

func (ips *IPSummary) GetTotalFiles() {
	ips.TotalFiles = len(ips.FileSummaries)
}

func (ips *IPSummary) GetStat() {
	var totalDS, err, suc int
	var totalPT time.Duration
	if len(ips.FileSummaries) == 0 {
		return
	}
	var min, max int
	var minPT, maxPT time.Duration
	if len(ips.FileSummaries) > 0 {
		ds := ips.FileSummaries[0].DownloadSize
		pt := ips.FileSummaries[0].ProcessTime
		minPT = pt
		maxPT = pt
		min = ds
		max = ds
	}

	for _, fs := range ips.FileSummaries {
		ds := fs.DownloadSize
		pt := fs.ProcessTime
		if fs.IsSucc() {
			totalDS += ds
			totalPT += pt
			suc += 1
		} else {
			err += 1
		}
		if ds < min {
			min = ds
		}
		if ds > max {
			max = ds
		}
		if pt < minPT {
			minPT = pt
		}
		if pt > maxPT {
			maxPT = pt
		}
	}
	ips.TotalProcessTime = totalPT
	if suc > 0 {
		ips.AvgProcessTime = totalPT / time.Duration(suc)
	}
	var sumsquares time.Duration

	for _, fs := range ips.FileSummaries {
		pt := fs.ProcessTime
		sumsquares += (pt - ips.AvgProcessTime) * (pt - ips.AvgProcessTime)
	}
	if suc > 0 {
		ips.StdDevProcessTime = time.Duration(math.Sqrt(
			float64(sumsquares / time.Duration(suc))))
	}
	ips.MaxProcessTime = maxPT
	ips.MinProcessTime = minPT
	ips.TotalDownloadErr = err
	ips.TotalDownloadSuc = suc
	ips.TotalDownloadSize = totalDS
	ips.MaxFileSize = max
	ips.MinFileSize = min
	if suc > 0 {
		ips.AvgFileSize = totalDS / suc
	}
	ips.TotalDownload = len(ips.FileSummaries)
}

type FileSummary struct {
	File         string
	DownloadSize int
	ProcessTime  time.Duration
	Err          *error
}

func (fs *FileSummary) IsSucc() bool {
	return fs.Err == nil
}

func NewFileSummary(file string, dSize int, pTime time.Duration, err *error) *FileSummary {
	return &FileSummary{
		File:         file,
		DownloadSize: dSize,
		ProcessTime:  pTime,
		Err:          err,
	}
}

type Summary struct {
	TotalFiles        int
	TotalDownloadSize int
	TotalDownload     int
	TotalDownloadErr  int
	TotalDownloadSuc  int
	MaxFileSize       int
	AvgFileSize       int
	MinFileSize       int
	DownloadErrs      map[string]int
	TotalProcessTime  time.Duration
	MaxProcessTime    time.Duration
	AvgProcessTime    time.Duration
	MinProcessTime    time.Duration
	StdDevProcessTime time.Duration
}
