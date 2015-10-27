package crawl

import (
	"bytes"
	"log"
	"sort"
	"testing"
)

func text2Segment(text []byte) (tline, segment []Typing) {
	base := 5
	lines := bytes.Split(text, []byte("\n"))
	lines = lines[1 : len(lines)-1]

	findLine := func(I int) int {
		for i := len(tline) - 1; i > -1; i-- {
			if tline[i].I == I {
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
			case '.':
				segment = append(segment, Typing{I: j})
			case '\\':
				k := findDownLine(i, j)
				if k < 0 {
					tline = append(tline, Typing{I: j, Type: DownTyping, High: base * 3, Low: base})
					k = len(tline) - 1
				}
				tline[k].I = j
				tline[k].Low = base
			case '/':
				k := findUpLine(i, j)
				if k < 0 {
					tline = append(tline, Typing{I: j, Type: UpTyping, High: base * 3, Low: base})
					k = len(tline) - 1
				}
				tline[k].Low = base
			}
		}
	}
	sort.Sort(TypingSlice(tline))
	for i := len(tline) - 1; i > -1; i-- {
		if tline[i].Type == DownTyping {
			tline[i].Price = tline[i].Low
		} else if tline[i].Type == UpTyping {
			tline[i].Price = tline[i].High
		}
	}
	sort.Sort(TypingSlice(segment))
	for i, c := 0, len(segment); i < c; i++ {
		if j := findLine(segment[i].I); j > -1 {
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
			log.Panicf("find segment[%d].I %d in tline fail, %s", i, segment[i].I, string(text))
		}
	}
	return
}

type test_text_data_pair struct {
	str     string
	line    []Typing
	segment []Typing
}

var tests_text_segment = []test_text_data_pair{
	{`
 /
/
      `,
		[]Typing{
			Typing{I: 1, Price: 25, Low: 5, High: 25, Type: UpTyping},
		}, nil,
	},
	{`
\
      `,
		[]Typing{
			Typing{I: 0, Price: 5, Low: 5, High: 15, Type: DownTyping},
		}, nil,
	},
	{`
\
 \/
      `,
		[]Typing{
			Typing{I: 1, Price: 5, Low: 5, High: 25, Type: DownTyping},
			Typing{I: 2, Price: 15, Low: 5, High: 15, Type: UpTyping},
		}, nil,
	},
	{`
    .    /
    /\  /
\  /  \/
 \/
 .
      `,
		[]Typing{
			Typing{I: 1, Price: 5, Low: 5, High: 25, Type: DownTyping},
			Typing{I: 4, Price: 35, Low: 5, High: 35, Type: UpTyping},
			Typing{I: 6, Price: 15, Low: 15, High: 35, Type: DownTyping},
			Typing{I: 9, Price: 45, Low: 15, High: 45, Type: UpTyping},
		},
		[]Typing{
			Typing{I: 1, Price: 5, Low: 5, High: 25, Type: BottomTyping},
			Typing{I: 4, Price: 35, Low: 5, High: 35, Type: TopTyping},
		},
	},
}

func test_line_i_price_type_equal(a, b []Typing) bool {
	if len(a) != len(b) {
		return false
	}
	for i, c := 0, len(a); i < c; i++ {
		if a[i].I != b[i].I || a[i].Type != b[i].Type || a[i].Price != b[i].Price {
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
