package crawl

import "github.com/golang/glog"

const (
	_ glog.Level = iota
	SegmentI
	LineI
	TypingI
	HttpI
	SegmentD
	LineD
	TypingD
	HttpD
	SegmentV
	LineV
	TypingV
	HttpV
)
