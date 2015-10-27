package crawl

import (
	"bytes"
	"testing"
)

func text2Segment(text []byte) (tline, segment []Typing) {
	base := 5
	lines := bytes.Split(text, []byte("\n"))
	lines = lines[1 : len(lines)-1]
	for _, l := range lines {
		if bytes.Contains(l, []byte("|")) {
			for i, c := 0, len(tline); i < c; i++ {
				if tline[i].High > 0 {
					tline[i].High = tline[i].High + base*2
					tline[i].Low = tline[i].Low + base*2
				}
			}
		}
		for i, c := range l {
			if len(tline) <= i {
				tline = append(tline, Typing{I: i})
			}

			switch c {
			case '^':
				tline[i].Type = UpTyping
			case 'v':
				tline[i].Type = DownTyping
			case '|':
				if tline[i].High == 0 {
					tline[i].High = base * 3
				}
				tline[i].Low = base
			}
		}
	}

	for i := len(tline) - 1; i > -1; i-- {
		if tline[i].Type == UpTyping {
			tline[i].Price = tline[i].High
		} else if tline[i].Type == DownTyping {
			tline[i].Price = tline[i].Low
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
^
|
|
      `,
		[]Typing{
			Typing{I: 0, Price: 25, Low: 5, High: 25, Type: UpTyping},
		}, nil,
	},
	{`
|
v
      `,
		[]Typing{
			Typing{I: 0, Price: 5, Low: 5, High: 15, Type: DownTyping},
		}, nil,
	},
	{`
|^
||
v
      `,
		[]Typing{
			Typing{I: 0, Price: 5, Low: 5, High: 25, Type: DownTyping},
			Typing{I: 1, Price: 15, Low: 5, High: 15, Type: UpTyping},
		}, nil,
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
