package crawl

import "time"

type Robot interface {
	Days_download(id string, start time.Time) ([]Tdata, error)
}
