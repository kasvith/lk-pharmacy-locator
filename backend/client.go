package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/go-resty/resty/v2"
)

func ProcessData(r io.Reader) ([]*PharmacyEntry, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}

	rows := f.GetRows("Sheet1")
	var data []*PharmacyEntry
	for _, row := range rows[1:] {
		data = append(data, &PharmacyEntry{
			District:       row[1],
			Area:           row[2],
			Name:           row[3],
			Address:        row[4],
			ContactNo:      strings.Fields(row[5]),
			PharmacistName: row[6],
			Owner:          row[7],
			WhatsApp:       strings.Fields(row[8]),
			Viber:          strings.Fields(row[9]),
			Email:          strings.Fields(row[10]),
		})
	}

	return data, nil
}

func FetchData() ([]*PharmacyEntry, error) {
	client := resty.New()
	resp, err := client.R().Get("https://docs.google.com/spreadsheets/d/1EzmE5KNIzy2cOE1OZdW7wo6MfLDmAq72relB9mxnbgo/export")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("error fetching data from google")
	}

	result, err := ProcessData(bytes.NewReader(resp.Body()))
	if err != nil {
		return nil, err
	}
	return result, nil
}
