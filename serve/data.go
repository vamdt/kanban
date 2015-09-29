package main

import (
	"encoding/json"
	"log"

	"./crawl"
)

var stocks crawl.Stocks

type data_response struct {
	stock_id string
	data     []byte
}

func (p *data_response) trans_data(stock *crawl.Stock, cb chan data_response) {
	data, err := json.Marshal(stock)
	if err != nil {
		log.Println(err)
	}
	p.data = data
	cb <- *p
}
