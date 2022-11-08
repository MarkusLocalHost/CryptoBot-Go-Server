package collectors

import (
	"archive/zip"
	"encoding/csv"
	"golang.org/x/text/encoding/charmap"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"
)

type Variant struct {
	FromID      string
	ToID        string
	MerchantID  string
	RateFrom    string
	RateTo      string
	Reserve     string
	BadReviews  string
	GoodReviews string
	Min         string
	Max         string
}

func GetExchangeRateVariants(from, to, fromType string) (map[string][]Variant, map[string]string) {
	fileinfo, err := os.Stat("../../tmp/bc_courses.zip")
	if err != nil {
		if os.IsNotExist(err) {
			downloadZipArchive()
		} else {
			log.Fatal(err)
		}
	}
	ctime := fileinfo.Sys().(*syscall.Stat_t).Ctim
	if time.Now().Sub(time.Unix(ctime.Sec, ctime.Nsec)).Seconds() > 45 {
		downloadZipArchive()
	}

	// Open zip file
	reader, err := zip.OpenReader("../../tmp/bc_courses.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	merchantIDs := make([]string, 0)
	merchants := make(map[string]string)
	result := make(map[string][]Variant)
	for _, f := range reader.File {
		// get all exchange variants
		if f.Name == "bm_rates.dat" {
			fc, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}
			defer func(fc io.ReadCloser) {
				err := fc.Close()
				if err != nil {
					panic(err)
				}
			}(fc)

			csvReader := csv.NewReader(charmap.Windows1251.NewDecoder().Reader(fc))
			records, err := csvReader.ReadAll()
			if err != nil {
				log.Fatal(err)
			}

			for _, record := range records {
				recordSlice := strings.Split(record[0], ";")
				variant := Variant{
					FromID:      recordSlice[0],
					ToID:        recordSlice[1],
					MerchantID:  recordSlice[2],
					RateFrom:    recordSlice[3],
					RateTo:      recordSlice[4],
					Reserve:     recordSlice[5],
					BadReviews:  strings.Split(recordSlice[6], ".")[0],
					GoodReviews: strings.Split(recordSlice[6], ".")[1],
					Min:         recordSlice[8],
					Max:         recordSlice[9],
				}

				if fromType == "bank" {
					banks := strings.Split(from, ";")
					if contains(banks, variant.FromID) && variant.ToID == to {
						merchantIDs = append(merchantIDs, variant.MerchantID)
						result[variant.FromID] = append(result[variant.FromID], variant)
					}
				} else if fromType == "wallet" {
					wallets := strings.Split(from, ";")
					if contains(wallets, variant.FromID) && variant.ToID == to {
						merchantIDs = append(merchantIDs, variant.MerchantID)
						result[variant.FromID] = append(result[variant.FromID], variant)
					}
				} else if fromType == "crypto" {
					to := strings.Split(to, ";")
					if variant.FromID == from && contains(to, variant.ToID) {
						merchantIDs = append(merchantIDs, variant.MerchantID)
						result[variant.FromID] = append(result[variant.FromID], variant)
					}
				}
			}
		}

		// get name merchant by id
		if f.Name == "bm_exch.dat" {
			fc, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}
			defer func(fc io.ReadCloser) {
				err := fc.Close()
				if err != nil {
					panic(err)
				}
			}(fc)

			csvReader := csv.NewReader(charmap.Windows1251.NewDecoder().Reader(fc))
			records, err := csvReader.ReadAll()
			if err != nil {
				log.Fatal(err)
			}

			for _, record := range records {
				recordSlice := strings.Split(record[0], ";")

				for _, merchantID := range merchantIDs {
					if merchantID == recordSlice[0] {
						merchants[merchantID] = recordSlice[1]
					}
				}
			}
		}
	}

	return result, merchants
}

func downloadZipArchive() {
	url := "http://api.bestchange.com/info.zip"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	// Create the file
	out, err := os.Create("../../tmp/bc_courses.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// Write to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
