package heater

import (
	"fmt"
	"net/http"
	"time"
)

type Task struct {
	schema      string
	hostname    string
	path        string
	file        string
	ip          string
	taskSummary TaskSummary
}

type TaskSummary struct {
	IP           string
	File         string
	ProcessTime  time.Duration
	DownloadSize int
	Err          *error
}

func NewTask(ip string, schema string, hostname string, path string, file string) *Task {
	return &Task{
		schema:   schema,
		hostname: hostname,
		path:     path,
		file:     file,
		ip:       ip,
		taskSummary: TaskSummary{
			IP:   ip,
			File: file,
		},
	}
}

func (t *Task) PutProcessTime(d time.Duration) {
	t.taskSummary.ProcessTime = d
}

func (t *Task) PutError(e *error) {
	t.taskSummary.Err = e
}

func (t *Task) PutResponse(r *http.Response) {
	return
	// for k, v := range r.Header {
	// 	fmt.Print(k)
	// 	fmt.Print(" : ")
	// 	fmt.Println(v)
	// }
}

func (t *Task) PutDownloadSize(s int) {
	t.taskSummary.DownloadSize = s
}

func (t *Task) GetIP() string {
	return t.ip
}

func (t *Task) GetPath() string {
	return t.path
}

func (t *Task) GetFile() string {
	return t.file
}

func (t *Task) GetHostname() string {
	return t.hostname
}

func (t *Task) GetURL() string {
	return fmt.Sprintf("%s://%s/%s/%s", t.schema, t.ip, t.path, t.file)
}

func (t *Task) GetSummary() TaskSummary {
	return t.taskSummary
}
