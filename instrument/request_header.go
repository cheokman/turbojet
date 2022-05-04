package instrument

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AliRequestHeader struct {
	HTTPVers               string    `json:"http_vers"`
	HTTPStatusCode         string    `json:"http_status_code"`
	Server                 string    `json:"server"`
	ContentType            string    `json:"content_type"`
	ContentLength          int32     `json:"content_length"`
	Connection             string    `json:"connection"`
	Date                   time.Time `json:"date"`
	LastModified           time.Time `json:"last_modified"`
	ETag                   string    `json:"etag"`
	AliSwiftGlobalSavetime int64     `json:"ali_swift_global_save_time"`
	Via                    string    `json:"via"`
	XCache                 string    `json:"x_cache"`
	XSwiftSaveTime         time.Time `json:"x_swift_save_time"`
	XSwiftCacheTime        int32     `json:"x_swift_cache_time"`
	TimingAllowOrigin      string    `json:"timing_allow_origin"`
	EagleId                string    `json:"eagle_id"`
	ContentEncoding        string    `json:"content_encoding"`
	TransferEncoding       string    `json:"transfer_encoding"`
	AcceptRanges           string    `json:"accept_ranges"`
	XSwiftError            string    `json:"x_swift_error"`
	Vary                   string    `json:"vary"`
	Age                    int32     `json:"age"`
}

type WangShuRequestHeader struct {
	HTTPVers       string
	HTTPStatusCode string
	Server         string
	ContentType    string
	ContentLength  int32
	Connection     string
	Date           time.Time
	LastModified   time.Time
	ETag           string
	XVia           string
	AcceptRanges   string
	Age            int32
}

func NewAliRequestHeader(data string) AliRequestHeader {
	var req AliRequestHeader
	var hv, hs string
	kvMap := make(map[string]string)
	dataArr := strings.Split(data, "\n")
	for i, item := range dataArr {
		if len(item) == 0 {
			continue
		}
		if i == 0 {
			versStatus := strings.Split(item, " ")
			if len(versStatus) == 3 {
				hv = versStatus[0]
				kvMap["HTTPVers"] = hv
				hs = versStatus[1] + " " + versStatus[2]
				kvMap["HTTPStatusCode"] = hs
			}
		} else {
			kv := strings.Split(item, ": ")
			if len(kv) == 1 {
				continue
			}
			kv[0] = strings.Replace(kv[0], "-", "", 10)
			kvMap[kv[0]] = kv[1]
		}

	}

	req = AliRequestHeader{
		HTTPVers:          kvMap["HTTPVers"],
		HTTPStatusCode:    kvMap["HTTPStatusCode"],
		Server:            kvMap["Server"],
		ContentType:       kvMap["ContentType"],
		TransferEncoding:  kvMap["TransferEncoding"],
		Connection:        kvMap["Connection"],
		Vary:              kvMap["Vary"],
		ETag:              kvMap["ETag"],
		Via:               kvMap["Via"],
		XCache:            kvMap["XCache"],
		TimingAllowOrigin: kvMap["TimingAllowOrigin"],
		EagleId:           kvMap["EagleId"],
		ContentEncoding:   kvMap["ContentEncoding"],
		AcceptRanges:      kvMap["AcceptRanges"],
		XSwiftError:       kvMap["XSwiftError"],
	}

	d, err := http.ParseTime(kvMap["Date"])
	if err == nil {
		req.Date = d
	}

	lm, err := http.ParseTime(kvMap["LastModified"])
	if err == nil {
		req.LastModified = lm
	}

	xSST, err := http.ParseTime(kvMap["XSwiftSaveTime"])
	if err == nil {
		req.XSwiftSaveTime = xSST
	}

	xSCT, err := strconv.ParseInt(kvMap["XSwiftCacheTime"], 10, 32)
	if err == nil {
		req.XSwiftCacheTime = int32(xSCT)
	}

	aSGS, err := strconv.ParseInt(kvMap["AliSwiftGlobalSavetime"], 10, 64)
	if err == nil {
		req.AliSwiftGlobalSavetime = aSGS
	}

	cl, err := strconv.ParseInt(kvMap["ContentLength"], 10, 32)
	if err == nil {
		req.ContentLength = int32(cl)
	}

	age, err := strconv.ParseInt(kvMap["Age"], 10, 32)
	if err == nil {
		req.ContentLength = int32(age)
	}
	return req
}
