package crawl

import (
	"bytes"
	"flag"
	"log"
	"sort"
	"strings"
	"testing"
)

var segment_files_flag = flag.String("segments", "", "the segment test files")

func text2Segment(text []byte) (tline, segment []Typing) {
	base := 5
	lines := bytes.Split(text, []byte("\n"))

	findLine := func(I int) int {
		for i := len(tline) - 1; i > -1; i-- {
			if tline[i].i == I {
				return i
			}
		}
		return -1
	}

	findUpLine := func(i, j int) int {
		index := -1
		for i, j = i-1, j+1; i > -1; i, j = i-1, j+1 {
			if len(lines[i]) > j && lines[i][j] == '/' {
				index = j
			} else {
				break
			}
		}
		return findLine(index)
	}

	findDownLine := func(i, j int) int {
		index := -1
		for i, j = i-1, j-1; i > -1 && j > -1; i, j = i-1, j-1 {
			if len(lines[i]) > j && lines[i][j] == '\\' {
				index = j
			} else {
				break
			}
		}
		return findLine(index)
	}

	for i, l := 0, len(lines); i < l; i++ {
		if bytes.IndexAny(lines[i], `/\`) > -1 {
			for j, c := 0, len(tline); j < c; j++ {
				tline[j].High = tline[j].High + base*2
				tline[j].Low = tline[j].Low + base*2
			}
		}
		for j, c := 0, len(lines[i]); j < c; j++ {
			switch lines[i][j] {
			case '*':
				segment = append(segment, Typing{i: j, Case1: true})
			case '?':
				segment = append(segment, Typing{i: j, Case1: false})
			case '\\':
				k := findDownLine(i, j)
				if k < 0 {
					tline = append(tline, downT(j, 0, base*3))
					k = len(tline) - 1
				}
				tline[k].Low = base
			case '/':
				k := findUpLine(i, j)
				if k < 0 {
					tline = append(tline, upT(j, 0, base*3))
					k = len(tline) - 1
				}
				tline[k].Low = base
			}
		}
	}
	sort.Sort(TypingSlice(tline))
	for i := len(tline) - 1; i > -1; i-- {
		tline[i].Time = tline[i].Time.UTC().AddDate(0, 0, i)
		tline[i].ETime = tline[i].Time
		if tline[i].Type == DownTyping {
			tline[i].Price = tline[i].Low
			tline[i].i = tline[i].i - 1 + (tline[i].High-tline[i].Low)/base/2
		} else if tline[i].Type == UpTyping {
			tline[i].Price = tline[i].High
		}
	}
	sort.Sort(TypingSlice(segment))
	for i, c := 0, len(segment); i < c; i++ {
		if j := findLine(segment[i].i); j > -1 {
			segment[i].i = j + 1
			segment[i].High = tline[j].High
			segment[i].Low = tline[j].Low
			segment[i].Price = tline[j].Price

			if tline[j].Type == DownTyping {
				segment[i].Type = BottomTyping
				if i > 0 {
					segment[i].High = segment[i-1].Price
				}
			} else if tline[j].Type == UpTyping {
				segment[i].Type = TopTyping
				if i > 0 {
					segment[i].Low = segment[i-1].Price
				}
			}
		} else {
			log.Panicf("find segment[%d].I %d in tline fail, %s", i, segment[i].i, string(text))
		}
	}
	return
}

type test_text_data_pair struct {
	str     string
	line    []Typing
	segment []Typing
}

func upT(i, low, high int) Typing {
	t := Typing{i: i, Type: UpTyping}
	t.Low = low
	t.High = high
	t.Price = high
	return t
}

func downT(i, low, high int) Typing {
	t := Typing{i: i, Type: DownTyping}
	t.Low = low
	t.High = high
	t.Price = low
	return t
}

func topT(i, low, high int, case1 bool) Typing {
	t := Typing{i: i, Type: TopTyping}
	t.Low = low
	t.High = high
	t.Price = high
	t.Case1 = case1
	return t
}

func bottomT(i, low, high int, case1 bool) Typing {
	t := Typing{i: i, Type: BottomTyping}
	t.Low = low
	t.High = high
	t.Price = low
	t.Case1 = case1
	return t
}

var tests_text_segment = []test_text_data_pair{
	{`
 /
/
      `,
		[]Typing{
			upT(1, 5, 25),
		}, nil,
	},
	{`
\
      `,
		[]Typing{
			downT(0, 5, 15),
		}, nil,
	},
	{`
\
 \/
      `,
		[]Typing{
			downT(1, 5, 25),
			upT(2, 5, 15),
		}, nil,
	},
	{`
         /\
        /  \
       /    \
  /\  /      \
 /  \/
/
      `,
		[]Typing{
			upT(2, 5, 35),
			downT(4, 15, 35),
			upT(9, 15, 65),
			downT(13, 25, 65),
		}, nil,
	},
	{`
    ?    /
    /\  /
\  /  \/
 \/
 *
      `,
		[]Typing{
			downT(1, 5, 25),
			upT(4, 5, 35),
			downT(6, 15, 35),
			upT(9, 15, 45),
		},
		[]Typing{
			bottomT(1, 5, 25, true),
			topT(2, 5, 35, false),
		},
	},
}

func test_line_i_price_type_equal(a, b []Typing) bool {
	if len(a) != len(b) {
		return false
	}
	for i, c := 0, len(a); i < c; i++ {
		if a[i].i != b[i].i || a[i].Type != b[i].Type || a[i].Price != b[i].Price {
			return false
		}
		if a[i].Case1 != b[i].Case1 {
			return false
		}
	}
	return true
}

func TestText2Segment(t *testing.T) {
	for i, pair := range tests_text_segment {
		lines, segments := text2Segment([]byte(pair.str))
		if !test_line_i_price_type_equal(lines, pair.line) {
			t.Error(
				"\nExample", i,
				"\nFor", pair.str,
				"\nexpected Line", pair.line,
				"\ngot", lines,
			)
		}
		if !test_line_i_price_type_equal(segments, pair.segment) {
			t.Error(
				"\nExample", i,
				"\nFor", pair.str,
				"\nexpected Segment", pair.segment,
				"\ngot", segments,
			)
		}
	}
}

func TestParseSegment(t *testing.T) {
	pattern := *segment_files_flag
	if len(pattern) < 1 {
		pattern = "**/*.segment"
	}
	tests_segments := load_test_desc_text_files(pattern)
	if tests_segments == nil {
		t.Fatal("load test files fail, pattern:", pattern)
	}
	t.Logf("load %d test files, pattern: %s", len(tests_segments), pattern)
	for i, d := range tests_segments {
		lines, segments := text2Segment([]byte(d.Text))
		td := Tdatas{}
		td.Typing.Line = lines
		td.ParseSegment()
		if !test_line_i_price_type_equal(segments, td.Segment.Data[1:]) {
			t.Error(
				"\nExample", i, d.File,
				"\nFor", d.Desc,
				"\nText", "\n"+strings.Replace(d.Text, " ", ".", -1),
				"\nexpected", segments,
				"\ngot", td.Segment.Data,
			)
		}
	}
}
