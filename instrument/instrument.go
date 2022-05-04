package instrument

import (
	"fmt"
	"net"
)

type Instrument struct {
	Location string  `json:"location"`
	ISP      string  `json:"isp"`
	Province string  `json:"province"`
	Points   []Point `json:"points"`
}

func NewInstrument(data []string) *Instrument {
	instr := &Instrument{
		Location: data[1],
		ISP:      data[2],
		Province: data[3],
	}
	point := NewPoint(data[4:14])
	request := NewAliRequestHeader(data[14])
	point.AliRequestHeader = request
	instr.Points = append(instr.Points, point)
	return instr
}

func (s *Instrument) AddPoint(data []string) {
	point := NewPoint(data[0:10])
	request := NewAliRequestHeader(data[10])
	point.AliRequestHeader = request
	s.Points = append(s.Points, point)
}

func (s *Instrument) GetIP() []string {
	var ips []string
	totalPoints := 0
	nilPoints := 0
	for _, p := range s.Points {
		totalPoints++
		included := false
		ipStr := fmt.Sprintf("%s", p.IP)
		for _, i := range ips {
			if i == ipStr {
				included = true
				break
			}
		}
		if included {
			continue
		}
		if net.ParseIP(ipStr) == nil {
			nilPoints++
			continue
		}
		ips = append(ips, ipStr)
	}
	if nilPoints > 0 {
		fmt.Printf("[%s] %d Nil IP Point over %d Total IP Point\n", s.Location, nilPoints, totalPoints)
	}
	return ips
}

func (s *Instrument) GetLocation() string {
	return s.Location
}

func (s *Instrument) GetISP() string {
	return s.ISP
}

func (s *Instrument) GetProvince() string {
	return s.Province
}
