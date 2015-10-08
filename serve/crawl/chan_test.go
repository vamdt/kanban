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
		false, false, false,
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
	Desc   string
	Data   []Tdata
	Typing []Typing
}

var tests_tdatas = []test_tdatas_pair{
	test_tdatas_pair{
		Desc: "Lesson 62 Fig 1, in the 3 k lines, the High of the k is the highest, and also the Low",
		Data: []Tdata{
			Tdata{High: 15, Low: 5},
			Tdata{High: 30, Low: 10},
			Tdata{High: 16, Low: 6},
		},
		Typing: []Typing{
			Typing{I: 1, Price: 30, Type: TopTyping},
		},
	},
	test_tdatas_pair{
		Desc: "Lesson 62 Fig 2, in the 3 k line, the High of the k is the lowest, and also the Low",
		Data: []Tdata{
			Tdata{High: 15, Low: 5},
			Tdata{High: 14, Low: 4},
			Tdata{High: 16, Low: 6},
		},
		Typing: []Typing{
			Typing{I: 1, Price: 4, Type: BottomTyping},
		},
	},
	test_tdatas_pair{
		Desc: "Lesson 62 Fig 3, the k line of the typing should not contain in another typing",
		Data: []Tdata{
			Tdata{High: 15, Low: 5},
			Tdata{High: 30, Low: 10},
			Tdata{High: 16, Low: 6},
			Tdata{High: 14, Low: 4},
			Tdata{High: 15, Low: 5},
		},
		Typing: []Typing{
			Typing{I: 1, Price: 30, Type: TopTyping},
		},
	},
	test_tdatas_pair{
		Desc: "Lesson 62 Fig 6, the contain-ship of TopTyping",
		Data: []Tdata{
			Tdata{High: 15, Low: 5},
			Tdata{High: 30, Low: 10},
			Tdata{High: 28, Low: 11},
			Tdata{High: 16, Low: 6},
		},
		Typing: []Typing{
      Typing{I: 2, Price: 30, Type: TopTyping},
		},
	},
	test_tdatas_pair{
		Desc: "Lesson 62 Fig 6, the contain-ship of BottomTyping",
		Data: []Tdata{
			Tdata{High: 20, Low: 10},
			Tdata{High: 14, Low: 4},
			Tdata{High: 13, Low: 5},
			Tdata{High: 18, Low: 9},
		},
		Typing: []Typing{
      Typing{I: 2, Price: 4, Type: BottomTyping},
		},
	},
	test_tdatas_pair{
		Desc: "Lesson 62 Fig 4, there should be 3 k lines between the TopTyping and the BottomTyping",
		Data: []Tdata{
			Tdata{High: 15, Low: 5},
			Tdata{High: 30, Low: 10},
			Tdata{High: 16, Low: 6},

			Tdata{High: 15, Low: 5},
			Tdata{High: 14, Low: 4},
			Tdata{High: 16, Low: 6},
		},
		Typing: []Typing{
			Typing{I: 1, Price: 30, Type: TopTyping},
		},
	},
	test_tdatas_pair{
		Desc: "Lesson 62 Fig 5, there should be 3 k lines between the TopTyping and the BottomTyping",
		Data: []Tdata{
			Tdata{High: 15, Low: 5},
			Tdata{High: 30, Low: 10},
			Tdata{High: 16, Low: 6},

			Tdata{High: 15, Low: 5},
			Tdata{High: 14, Low: 4},
			Tdata{High: 13, Low: 3},
			Tdata{High: 17, Low: 7},
		},
		Typing: []Typing{
			Typing{I: 1, Price: 30, Type: TopTyping},
			Typing{I: 5, Price: 3, Type: BottomTyping},
		},
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
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestParseTyping(t *testing.T) {
	for i, d := range tests_tdatas {
		td := Tdatas{Data: d.Data}
		td.ParseTyping()
		if !test_is_typing_equal(t, d.Typing, td.Typing) {
			t.Error(
        "Test", i,
				"For", d.Desc,
				"expected", d.Typing,
				"got", td.Typing,
			)
		}
	}
}
