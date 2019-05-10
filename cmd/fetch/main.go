package main

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/nkozyra/ipcarta"
)

const (
	countryCSVURL = "https://geolite.maxmind.com/download/geoip/database/GeoLite2-Country-CSV.zip"
	cityCSVURL    = "https://geolite.maxmind.com/download/geoip/database/GeoLite2-City-CSV.zip"
)

var (
	lang          string
	ipLookups     map[string]ipcarta.Network
	ips           []ipcarta.Network
	elasticsearch string
)

func init() {
	fmt.Println("IP Carta")
	flag.StringVar(&lang, "lang", "en", "Language to use")
	flag.StringVar(&elasticsearch, "el", "", "Elastic search address")
	flag.Parse()
}

func main() {
	ipLookups = make(map[string]ipcarta.Network)
	ipcarta.Init(ipcarta.Config{
		ElasticSearchHost: elasticsearch},
	)
	fetchCountries()
	setIPs()
}

func readZip(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func setIPs() {
	for _, v := range ips {
		ipcarta.Set(v)
	}
}

func fetchCountries() {
	coResp, err := http.Get(countryCSVURL)
	if err != nil {
		panic(err)
	}
	defer coResp.Body.Close()
	coData, err := ioutil.ReadAll(coResp.Body)
	if err != nil {
		panic(err)
	}
	var countryData []byte
	zr, _ := zip.NewReader(bytes.NewReader(coData), int64(len(coData)))
	for _, zf := range zr.File {
		fn := strings.Split(zf.Name, "/")[1]
		if fn == "GeoLite2-Country-Blocks-IPv4.csv" {
			data, _ := readZip(zf)
			out, err := os.Create("tmp/countries.csv")
			if err != nil {
				panic(err)
			}
			defer out.Close()
			out.Write(data)
			out.Close()
			countryData = data
		}
		if fn == fmt.Sprintf("GeoLite2-Country-Locations-%s.csv", lang) {
			data, _ := readZip(zf)
			out, err := os.Create("tmp/country-lookups.csv")
			if err != nil {
				panic(err)
			}
			defer out.Close()
			out.Write(data)
			out.Close()
			r := csv.NewReader(bytes.NewReader(data))
			for {
				record, err := r.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatal(err)
				}
				geoID := record[0]
				continentCode := record[2]
				continentName := record[3]
				countryCode := record[4]
				countryName := record[5]

				if _, ok := ipLookups[geoID]; !ok {
					ipLookups[geoID] = ipcarta.Network{
						Network:        geoID,
						ContinentCode:  continentCode,
						ContinentName:  continentName,
						CountryISOCode: countryCode,
						CountryName:    countryName,
					}
				}
			}
		}
	}
	r := csv.NewReader(bytes.NewReader(countryData))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		geoID := record[1]
		if _, ok := ipLookups[geoID]; ok {
			ip := ipLookups[geoID]
			ip.Network = record[0]
			ip.IsAnonymousProxy = record[4] == "1"
			ips = append(ips, ip)
		}
	}
}
