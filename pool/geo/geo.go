package geo

import (
	"io/ioutil"
	"net/http"
)

type Geo struct {
	Status int `json:"status"`
	Exact float64 `json:"exact"`
	Gps bool `json:"gps"`
	Gcj Gcj `json:"gcj"`
	Bd09 Bd09 `json:"bd09"`
	URL string `json:"url"`
	Address string `json:"address"`
	IP string `json:"ip"`
	Cache int `json:"cache"`
}
type Gcj struct {
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
}
type Bd09 struct {
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
}

func GetGeo() (string,error) {
	resp, err := http.Get("https://api.asilu.com/geo/")
	if err != nil {
		return "",err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body),err
}
