package main

import (
	"encoding/json"
	"log"

	"./crawl"
)

type data_response struct {
	stock_id string
	data     []byte
}

func load_data(stock_id string, cb chan data_response) {
	res := data_response{stock_id: stock_id}
	stock := crawl.Stock{Id: stock_id}
	new_num := stock.Days_sync(db)
	log.Println("new sync", new_num)
	data, err := json.Marshal(stock)
	if err != nil {
		log.Println(err)
	}
	res.data = data
	cb <- res
}
