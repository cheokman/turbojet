package instrument

import (
	"bytes"
	"encoding/json"
	"turbojet/cli"

	"github.com/PuerkitoBio/goquery"
)

type InstrumentSlice struct {
	Instruments []*Instrument `json:"instruments"`
}

func (s *InstrumentSlice) GetLocations() []string {
	var locations []string
	for _, l := range s.Instruments {
		lt := l.GetLocation()
		contained := false
		for _, el := range locations {
			if lt == el {
				contained = true
				break
			}
		}
		if !contained {
			locations = append(locations)
		}
	}
	return locations
}

func (s *InstrumentSlice) GetProvinces() []string {
	var provinces []string
	for _, p := range s.Instruments {
		pv := p.GetProvince()
		contained := false
		for _, ep := range provinces {
			if ep == pv {
				contained = true
				break
			}
		}
		if !contained {
			provinces = append(provinces, pv)
		}
	}
	return provinces
}

func (s *InstrumentSlice) GetIPs() []string {
	var ips []string
	for _, i := range s.Instruments {
		_ips := i.GetIP()
		for _, j := range _ips {
			included := false
			for _, ip := range ips {
				if ip == j {
					included = true
					break
				}
			}
			if included {
				continue
			}
			ips = append(ips, j)
		}
	}
	return ips
}

func Parse(c *cli.Context, source string) (InstrumentSlice, error) {
	var currentInstrument *Instrument
	var instrSlice InstrumentSlice //[]*Instrument
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(source)))
	w := c.Writer()
	if err != nil {
		cli.Errorf(w, "Instrument source error: %s\n", err)
		return instrSlice, err
	}
	rows := getTableRows(doc)
	for i, r := range rows {
		length := len(r)
		if length == 0 {
			continue
		}
		if i == 1 {
			continue
		}
		if length == 11 {
			currentInstrument.AddPoint(r)
		}
		if length == 16 {
			currentInstrument = NewInstrument(r)
			instrSlice.Instruments = append(instrSlice.Instruments, currentInstrument)
		}
	}
	return instrSlice, nil
}

func LoadCache(c *cli.Context, data []byte) (InstrumentSlice, error) {
	var instruments InstrumentSlice
	err := json.Unmarshal(data, &instruments)
	if err != nil {
		cli.Errorf(c.Writer(), "Instrument load from cache error: %s\n", err)
		return instruments, err
	}
	return instruments, nil
}

func getTableRows(doc *goquery.Document) [][]string {
	var headings, row []string
	var rows [][]string
	doc.Find("table").Each(func(index int, tablehtml *goquery.Selection) {
		tablehtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
			rowhtml.Find("th").Each(func(indexth int, tableheading *goquery.Selection) {
				headings = append(headings, tableheading.Text())
			})
			rowhtml.Find("td").Each(func(indexth int, tablecell *goquery.Selection) {
				row = append(row, tablecell.Text())
			})
			rows = append(rows, row)
			row = nil
		})
	})
	return rows
}
