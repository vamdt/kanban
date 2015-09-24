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

func (p *data_response) trans_data(stock *crawl.Stock, cb chan data_response) {
	data, err := json.Marshal(stock)
	if err != nil {
		log.Println(err)
	}
	p.data = data
	cb <- *p
}

func load_data(stock_id string, cb chan data_response) {
	res := data_response{stock_id: stock_id}
	stock := crawl.Stock{Id: stock_id}

	new_num := stock.Days_update(db)
	res.trans_data(&stock, cb)
	log.Println("new days update", new_num)

	new_num = stock.M30s_update(db)
	res.trans_data(&stock, cb)
	log.Println("new m30 update", new_num)

	new_num = stock.M5s_update(db)
	res.trans_data(&stock, cb)
	log.Println("new m5 update", new_num)

	new_num = stock.Ticks_update(db)
	res.trans_data(&stock, cb)
	log.Println("new ticks update", new_num)
}
