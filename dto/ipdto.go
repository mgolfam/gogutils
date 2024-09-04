package dto

import "strconv"

type LangConf struct {
	Rtl  bool
	Lang string
}

func (l *LangConf) SetRtl(rtl string) {
	b, err := strconv.ParseBool(rtl)
	if err != nil {
		return
	}

	l.Rtl = b
}

func (l *LangConf) Default() {
	// l.Lang = "fa"
	// l.Rtl = true

	l.Lang = "en"
	l.Rtl = false
}

type IpInfo struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}
