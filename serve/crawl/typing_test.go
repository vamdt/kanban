package crawl

import (
	"bytes"
	"flag"
	"sort"
	"testing"
)

var line_files_flag = flag.String("lines", "", "the line test files")
var typing_files_flag = flag.String("typings", "", "the typing test files")

func text2Tdatas(text []byte) Tdatas {
	tds := Tdatas{}
	base := 5
	td := []Tdata{}
	typing := []Typing{}
	tline := []Typing{}
	lines := bytes.Split(text, []byte("\n"))
	for _, l := range lines {
    if bytes.Count(l, []byte(" ")) == len(l) {
      continue
    }
		if bytes.IndexAny(l, "|-_") > -1 {
			for i, c := 0, len(td); i < c; i++ {
				if td[i].High > 0 {
					td[i].High = td[i].High + base*2
					td[i].Low = td[i].Low + base*2
				}
			}
		}
		for i, c := range l {
			if len(td) <= i {
				td = append(td, Tdata{})
			}

			switch c {
			case 'L':
				tline = append(tline, Typing{I: i, Type: TopTyping})
			case 'l':
				tline = append(tline, Typing{I: i, Type: BottomTyping})
			case '^':
				typing = append(typing, Typing{I: i, Type: TopTyping})
			case '.':
				fallthrough
			case 'v':
				typing = append(typing, Typing{I: i, Type: BottomTyping})
			case '|':
				if td[i].High == 0 {
					td[i].High = base * 3
				}
				td[i].Low = base
			case '-':
				if td[i].High == 0 {
					td[i].High = base * 2
				}
				td[i].Low = base * 2
			case '_':
				if td[i].High == 0 {
					td[i].High = base
				}
				td[i].Low = base
			}
		}
	}
	tds.Data = td
	for i, c := 0, len(typing); i < c; i++ {
		typing[i].High = td[typing[i].I].High
		typing[i].Low = td[typing[i].I].Low
		if typing[i].Type == TopTyping {
			typing[i].Price = td[typing[i].I].High
		} else if typing[i].Type == BottomTyping {
			typing[i].Price = td[typing[i].I].Low
		}
	}
	sort.Sort(TypingSlice(typing))
	tds.Typing.Data = typing

	if llen := len(tline); llen > 0 {
		sort.Sort(TypingSlice(tline))
		for i := llen - 1; i > -1; i-- {
			tline[i].High = td[tline[i].I].High
			tline[i].Low = td[tline[i].I].Low
			if tline[i].Type == TopTyping {
				tline[i].Price = tline[i].High
				tline[i].Type = DownTyping
			} else if tline[i].Type == BottomTyping {
				tline[i].Price = tline[i].Low
				tline[i].Type = UpTyping
			}
		}
		tds.Typing.Line = tline[:llen-1]
	}
	return tds
}

type test_text_tdata_pair struct {
	str        string
	exp_td     []Tdata
	exp_typing []Typing
	exp_line   []Typing
}

var tests_text_tdata = []test_text_tdata_pair{
	{`
|
      `,
		[]Tdata{
			Tdata{High: 15, Low: 5},
		}, nil, nil,
	},
	{`
|
|
      `,
		[]Tdata{
			Tdata{High: 25, Low: 5},
		}, nil, nil,
	},
	{`
^
|
      `,
		[]Tdata{
			Tdata{High: 15, Low: 5},
		},
		[]Typing{
			Typing{I: 0, Price: 15, Type: TopTyping},
		}, nil,
	},
	{`
|^
||
      `,
		[]Tdata{
			Tdata{High: 25, Low: 5},
			Tdata{High: 15, Low: 5},
		},
		[]Typing{
			Typing{I: 1, Price: 15, Type: TopTyping},
		}, nil,
	},
	{`
|
.
      `,
		[]Tdata{
			Tdata{High: 15, Low: 5},
		},
		[]Typing{
			Typing{I: 0, Price: 5, Type: BottomTyping},
		}, nil,
	},
	{`
    L
    ^
    |
    | |
 |  | |_|
||-_||| |||
 |  | | | |
      |
      .
      l
      `,
		[]Tdata{
			Tdata{High: 35, Low: 25},
			Tdata{High: 45, Low: 15},
			Tdata{High: 30, Low: 30},
			Tdata{High: 25, Low: 25},
			Tdata{High: 65, Low: 15},

			Tdata{High: 35, Low: 25},
			Tdata{High: 55, Low: 5},
			Tdata{High: 35, Low: 35},
			Tdata{High: 45, Low: 15},
			Tdata{High: 35, Low: 25},

			Tdata{High: 35, Low: 15},
		},
		[]Typing{
			Typing{I: 4, Price: 65, Type: TopTyping},
			Typing{I: 6, Price: 5, Type: BottomTyping},
		},
		[]Typing{
			Typing{I: 4, Price: 65, Type: DownTyping},
		},
	},
}

func test_tdata_high_low_equal(a, b []Tdata) bool {
	if len(a) != len(b) {
		return false
	}
	for i, c := 0, len(a); i < c; i++ {
		if a[i].High != b[i].High || a[i].Low != b[i].Low {
			return false
		}
	}
	return true
}

func test_typing_i_price_type_equal(a, b []Typing) bool {
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

func TestText2Tdata(t *testing.T) {
	for i, pair := range tests_text_tdata {
		tds := text2Tdatas([]byte(pair.str))
		if !test_tdata_high_low_equal(tds.Data, pair.exp_td) {
			t.Error(
				"\nExample", i,
				"\nFor", pair.str,
				"\nexpected Tdata", pair.exp_td,
				"\ngot", tds.Data,
			)
		}
		if !test_typing_i_price_type_equal(tds.Typing.Data, pair.exp_typing) {
			t.Error(
				"\nExample", i,
				"\nFor", pair.str,
				"\nexpected Typing", pair.exp_typing,
				"\ngot", tds.Typing,
			)
		}
		if !test_typing_i_price_type_equal(tds.Typing.Line, pair.exp_line) {
			t.Error(
				"\nExample", i,
				"\nFor", pair.str,
				"\nexpected Line", pair.exp_line,
				"\ngot", tds.Typing.Line,
			)
		}
	}
}

type test_typing_pair struct {
	tdata      [3]Tdata
	is_top     bool
	is_bottom  bool
	is_contain bool
}

var tests_typing = []test_typing_pair{
	test_typing_pair{
		[3]Tdata{
			Tdata{High: 100, Low: 90},
			Tdata{High: 200, Low: 100},
			Tdata{High: 150, Low: 80},
		},
		true, false, false,
	},
	test_typing_pair{
		[3]Tdata{
			Tdata{High: 100, Low: 90},
			Tdata{High: 100, Low: 100},
			Tdata{High: 150, Low: 80},
		},
		false, false, true,
	},
	test_typing_pair{
		[3]Tdata{
			Tdata{High: 100, Low: 90},
			Tdata{High: 200, Low: 90},
			Tdata{High: 150, Low: 80},
		},
		false, false, true,
	},
	test_typing_pair{
		[3]Tdata{
			Tdata{High: 100, Low: 90},
			Tdata{High: 200, Low: 70},
			Tdata{High: 150, Low: 80},
		},
		false, false, true,
	},
	test_typing_pair{
		[3]Tdata{
			Tdata{High: 100, Low: 90},
			Tdata{High: 90, Low: 70},
			Tdata{High: 150, Low: 80},
		},
		false, true, false,
	},
	test_typing_pair{
		[3]Tdata{
			Tdata{High: 200, Low: 90},
			Tdata{High: 140, Low: 100},
			Tdata{High: 150, Low: 80},
		},
		false, false, true,
	},
}

func TestIsTopTyping(t *testing.T) {
	for i, td := range tests_typing {
		if td.is_top != IsTopTyping(&td.tdata[0], &td.tdata[1], &td.tdata[2]) {
			t.Error(
				"Test", i,
				"For", td.tdata,
				"expected", td.is_top,
				"got", !td.is_top,
			)
		}
	}
}

func TestIsBottomTyping(t *testing.T) {
	for i, td := range tests_typing {
		if td.is_bottom != IsBottomTyping(&td.tdata[0], &td.tdata[1], &td.tdata[2]) {
			t.Error(
				"Test", i,
				"For", td.tdata,
				"expected", td.is_bottom,
				"got", !td.is_bottom,
			)
		}
	}
}

func TestContain(t *testing.T) {
	for i, td := range tests_typing {
		if td.is_contain != Contain(&td.tdata[0], &td.tdata[1]) {
			t.Error(
				"Test", i,
				"For", td.tdata,
				"expected", td.is_contain,
				"got", !td.is_contain,
			)
		}
	}
}

func test_is_typing_equal(t *testing.T, a, b []Typing) bool {
	if len(a) != len(b) {
		return false
	}
	if len(a) == 0 {
		return true
	}
	for i := 0; i < len(a); i++ {
		if a[i].I != b[i].I || a[i].Price != b[i].Price || a[i].Type != b[i].Type {
			return false
		}
	}
	return true
}

func TestParseTyping(t *testing.T) {
	pattern := *typing_files_flag
	if len(pattern) < 1 {
		pattern = "**/*.typing"
	}
	tests_tdatas := load_test_desc_text_files(pattern)
	if tests_tdatas == nil {
		t.Fatal("load test files fail, pattern:", pattern)
	}
	t.Logf("load %d test files, pattern: %s", len(tests_tdatas), pattern)
	for i, d := range tests_tdatas {
		exp := text2Tdatas([]byte(d.Text))
		td := Tdatas{Data: exp.Data}
		td.ParseTyping()
		if !test_is_typing_equal(t, exp.Typing.Data, td.Typing.Data) {
			t.Error(
				"\nExample", i,
				"\nFor", d.Desc,
				"\nText", "\n"+d.Text,
				"\nexpected", exp.Typing,
				"\ngot", td.Typing,
			)
		}
	}
}

func TestLinkTyping(t *testing.T) {
	pattern := *line_files_flag
	if len(pattern) < 1 {
		pattern = "**/*.line"
	}
	tests_lines := load_test_desc_text_files(pattern)
	if tests_lines == nil {
		t.Fatal("load test files fail, pattern:", pattern)
	}
	t.Logf("load %d test files, pattern: %s", len(tests_lines), pattern)
	for i, d := range tests_lines {
		exp := text2Tdatas([]byte(d.Text))
		td := Tdatas{Data: exp.Data}
		td.ParseTyping()
		td.Typing.LinkTyping()
		if !test_is_typing_equal(t, exp.Typing.Line, td.Typing.Line) {
			t.Error(
				"\nExample", i,
				"\nFor", d.Desc,
				"\nText", "\n"+d.Text,
				"\nexpected", exp.Typing.Line,
				"\ngot", td.Typing.Line,
			)
		}
	}
}
