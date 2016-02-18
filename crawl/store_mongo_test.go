// +build mongo

package crawl

import (
	"encoding/json"
	"testing"
)

func TestBsonMarshal(t *testing.T) {
	type A struct {
		High, Low int
	}
	type B struct {
		A `bson:",inline"`
	}
	b := B{}
	b.High = 1
	b.Low = 1
	m, err := data2BsonM(b)
	if err != nil {
		t.Error("For", "bson marchal fail", err)
	}
	buf, _ := json.Marshal(m)
	s := string(buf)
	exp := `{"high":1,"low":1}`
	if s != exp {
		t.Error(
			"For", "bson marchal",
			"expected", exp,
			"got", s,
		)
	}
}
