package crawl

import "testing"

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
		false, false, false,
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

type test_tdatas_pair struct {
	Desc string
	Text string
}

var tests_tdatas = []test_tdatas_pair{
	test_tdatas_pair{
		Desc: "Lesson 62 Fig 1, in the 3 k lines, the High of the k is the highest, and also the Low",
		Text: `
 ^
 |
 ||
|||
| |
|
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 62 Fig 2, in the 3 k line, the Low of the k is the lowest, and also the High",
		Text: `
|
| |
|||
 ||
 |
 .
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 62 Fig 3, the k line of the typing should not contain in another typing",
		Text: `
 ^
 |
 |
 ||
||| |
| |||
|  ||
   |
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 62 Fig 6, the contain-ship of TopTyping",
		Text: `
 ^
 |
 ||
 |||
|| |
|  ||
|   |
    |
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 62 Fig 6, the contain-ship of BottomTyping",
		Text: `
|
|  |
|| |
 |||
 ||
 |
 .
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 62 Fig 4, there should be 3 k lines between the TopTyping and the BottomTyping",
		Text: `
 ^
 |
 ||
|||| |
| ||||
|  |||
|   |
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 62 Fig 5, there should be 3 k lines between the TopTyping and the BottomTyping",
		Text: `
 ^
 |
 ||
||||  |
| |||||
|  ||||
|   ||
     |
     .
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 62 Fig 5, should not have two TopTyping together",
		Text: `
         ^
     |   |
     ||  ||
    ||| |||||
    ||||| |||||       |
   ||||||  ||||||    ||
   |||||     | |||||||
 ||| | |        ||||||
|||  | |         |
||               .
|
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 62 Study Case 9, Top should not lower then Bottom",
		Text: `
   |
| ||
||||
||||  |
|||| |||
||| ||||
 |  || |
 |  |
 .
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 65 Fig 4, case TopTop, should skip the first Top",
		Text: `
       ^
       |
      |||
   | || |
  ||||
 || |
||
|
    `,
	},
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
	for i, d := range tests_tdatas {
		exp := Text2Tdatas([]byte(d.Text))
		td := Tdatas{Data: exp.Data}
		td.ParseTyping()
		if !test_is_typing_equal(t, exp.Typing, td.Typing) {
			t.Error(
				"\nExample", i,
				"\nFor", d.Desc,
				"\nText", d.Text,
				"\nexpected", exp.Typing,
				"\ngot", td.Typing,
			)
		}
	}
}
