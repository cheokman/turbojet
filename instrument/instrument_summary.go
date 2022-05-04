package instrument

import "time"

type InstrumentSummary struct {
	TotalInstrument      int           `json:"total_instrument"`
	TotalISP             int           `json:"total_isp"`
	TotalProvince        int           `json:"total_provice"`
	UnitIPs              int           `json:"unit_ips"`
	UnitLocation         int           `json:"unit_location"`
	Successful           bool          `json:"successful"`
	AvgTotalTime         time.Duration `json:"avg_total_time"`
	AvgResolvedTime      time.Duration `json:"avg_resolved_time"`
	AvgConnectTime       time.Duration `json:"avg_connect_time"`
	AvgDownloadTime      time.Duration `json:"avg_download_time"`
	SameFileSize         bool          `json:"same_file_size"`
	SameFileDownloadSize bool          `json:"same_file_download_size"`
	AvgDownloadSpeedMBS  float32       `json:"agv_download_speed_mbs"`
}
