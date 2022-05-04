package instrument

import (
	"net"
	"strconv"
	"strings"
	"time"
)

type Point struct {
	IP                   net.IP        `json:"ip"`
	IPLocation           string        `json:"ip_location"`
	HTTPStatus           string        `json:"http_status"`
	TotalTime            time.Duration `json:"total_time"`
	ResolvedTime         time.Duration `json:"resovled_time"`
	ConnectTime          time.Duration `json:"connect_time"`
	DownloadTime         time.Duration `json:"download_time"`
	ContentLength        int32         `json:"content_length"`
	DownloadLength       int32         `json:"download_length"`
	DownloadSpeed        float32       `json:"download_speed"`
	AliRequestHeader     `json:"ali_request_header"`
	WangShuRequestHeader `json:"wangshu_request_header"`
}

func NewPoint(data []string) Point {
	var p Point

	p = Point{
		IP:             net.ParseIP(data[0]),
		IPLocation:     data[1],
		HTTPStatus:     data[2],
		TotalTime:      parseDuration(data[3]),
		ResolvedTime:   parseDuration(data[4]),
		ConnectTime:    parseDuration(data[5]),
		DownloadTime:   parseDuration(data[6]),
		ContentLength:  parseHeaderLength(data[7]),
		DownloadLength: parseHeaderLength(data[8]),
		DownloadSpeed:  parseSpeed(data[9]),
	}
	return p
}

func parseDuration(d string) time.Duration {

	d = strings.Replace(d, "s", "", 1)
	duration, err := strconv.ParseFloat(d, 32)
	if err != nil {
		return time.Duration(0 * time.Second)
	}
	return time.Duration(int32(duration * float64(time.Second)))
}

func parseHeaderLength(l string) int32 {
	d := strings.Replace(l, "KB", "", 1)
	length, err := strconv.ParseInt(d, 10, 32)
	if err != nil {
		return 0
	}

	return int32(length)
}

func parseSpeed(s string) float32 {
	d := strings.Replace(s, "MB/s", "", 1)
	speed, err := strconv.ParseFloat(d, 32)
	if err != nil {
		return 0.0
	}

	return float32(speed)
}
