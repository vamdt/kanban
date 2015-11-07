package crawl

import (
	"bytes"
	"log"
	"sort"
	"strings"
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
			case '*':
				segment = append(segment, Typing{I: j})
			case '\\':
				k := findDownLine(i, j)
				if k < 0 {
					tline = append(tline, Typing{I: j, Type: DownTyping, High: base * 3})
					k = len(tline) - 1
				}
				tline[k].Low = base
			case '/':
				k := findUpLine(i, j)
				if k < 0 {
					tline = append(tline, Typing{I: j, Type: UpTyping, High: base * 3})
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
			tline[i].I = tline[i].I - 1 + (tline[i].High-tline[i].Low)/base/2
		} else if tline[i].Type == UpTyping {
			tline[i].Price = tline[i].High
		}
	}
	sort.Sort(TypingSlice(segment))
	for i, c := 0, len(segment); i < c; i++ {
		if j := findLine(segment[i].I); j > -1 {
			segment[i].I = j + 1
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
         /\
        /  \
       /    \
  /\  /      \
 /  \/
/
      `,
		[]Typing{
			Typing{I: 2, Price: 35, Low: 5, High: 35, Type: UpTyping},
			Typing{I: 4, Price: 15, Low: 15, High: 35, Type: DownTyping},
			Typing{I: 9, Price: 65, Low: 15, High: 65, Type: UpTyping},
			Typing{I: 13, Price: 25, Low: 25, High: 65, Type: DownTyping},
		}, nil,
	},
	{`
    *    /
    /\  /
\  /  \/
 \/
 *
      `,
		[]Typing{
			Typing{I: 1, Price: 5, Low: 5, High: 25, Type: DownTyping},
			Typing{I: 4, Price: 35, Low: 5, High: 35, Type: UpTyping},
			Typing{I: 6, Price: 15, Low: 15, High: 35, Type: DownTyping},
			Typing{I: 9, Price: 45, Low: 15, High: 45, Type: UpTyping},
		},
		[]Typing{
			Typing{I: 1, Price: 5, Low: 5, High: 25, Type: BottomTyping},
			Typing{I: 2, Price: 35, Low: 5, High: 35, Type: TopTyping},
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

var tests_segments = []test_tdatas_pair{
	test_tdatas_pair{
		Desc: "Lesson 67 Study Fig 1, Case 1 standard",
		Text: `
        *
        /\
       /  \  /\
  /\  /    \/  \
 /  \/          \
/
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 67 Study Fig 2, Case 1 standard extend",
		Text: `
         *
         /\
        /  \                /
       /    \  /\          /
  /\  /      \/  \    /\  /
 /  \/            \  /  \/
/                  \/
                   *
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 67 Study Case 2 Fig 1, Case 2 standard",
		Text: `
            g2
           /\      g3
          /  \    /\
         /    \  /  \
        /      \/    \
   g1  /        d3    \
  /\  /                \
 /  \/                  d4
/    d2
d1
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 67 Study Case 2 Fig 2, Case 2 standard ensure",
		Text: `
            *g2
            /\                      g5
           /  \      g3            /
          /    \    /\            /
         /      \  /  \      g4  /
        /        \/    \    /\  /
   g1  /          d3    \  /  \/
  /\  /                  \/    d5
 /  \/                   *d4
/    d2
d1
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 67 Study Case 2 Fig 3, A(Case 2) + B + C fail",
		Text: `
                                      /g5
             g2                      /
            /\                      /
           /  \      g3            /
          /    \    /\            /
         /      \  /  \      g4  /
        /        \/    \    /\  /
   g1  /          d3    \  /  \/
  /\  /                  \/    d5
 /  \/                   d4
/    d2
d1
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 67 Study Case 8 special",
		Text: `
              *
              /\
             /  \
            /    \  /\
   /\      /      \/  \
  /  \/\  /            \
 /      \/              \
/
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 71 Study Case 2 - 1",
		Text: `
                         *8
                   6     /\
                  /\    /  \/\10
         4       /  \  /    9 \
    2   /\      /    \/7       11
   /\  /  \    /
  /  \/    \  /
 /    3     \/
/1           5
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 71 Study Case 2 - 2",
		Text: `
             4
            /\      6
           /  \    /\        8
          /    \  /  \      /\
         /      \/    \    /  \
    2   /        5     \  /    \
   /\  /                \/      9
  /  \/                  7
 /    3
/1
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 71 Study Case 3",
		Text: `
                          8
                   6     /\
                  /\    /  \
         4       /  \  /    \
    2   /\      /    \/7     \
   /\  /  \    /              \9
  /  \/    \  /
 /    3     \/
/1           5
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 71 Study Case 4 - 1",
		Text: `
                          /\8
                   6     /  \        /\10
                  /\    /    \      /  \
         4       /  \  /      \    /    \
    2   /\      /    \/        \  /      11
   /\  /  \    /      7         \/
  /  \/    \  /                  9
 /    3     \/
/1           5
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 71 Study Case 4 - 2",
		Text: `
                                             /12
                          /\8               /
                   6     /  \        /\ 10 /
                  /\    /    \      /  \  /
         4       /  \  /      \    /    \/
    2   /\      /    \/        \  /      11
   /\  /  \    /      7         \/
  /  \/    \  /                  9
 /    3     \/
/1           5
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 71 Study Case 4 - 3",
		Text: `
                          *8
                          /\          10
                   6     /  \        /\
                  /\    /    \      /  \  /\  12
         4       /  \  /      \    /    \/  \
    2   /\      /    \/        \  /      11  \
   /\  /  \    /      7         \/            \
  /  \/    \  /                  9             \
 /    3     \/                                  13
/1           5
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 71 Study Case 5, as 4-3",
		Text: `
                          *8
                          /\          10
                   6     /  \        /\
                  /\    /    \      /  \  /\  12
         4       /  \  /      \    /    \/  \
    2   /\      /    \/        \  /      11  \
   /\  /  \    /      7         \/            \
  /  \/    \  /                  9             \
 /    3     \/                                  13
/1           5
    `,
	},
	test_tdatas_pair{
		Desc: `Lesson 77 Case 81-82 & Lesson 78 // Fuzzy 请参考书中原图`,
		Text: `
            *80
            /\                    b                    82
           /  \                  /\                    h
          /    \                /  \              f   /
         /      \    /\        /    \      d     /\  /
        /        \  /  \      /      \    /\    /  \/
   /\  /          \/    \    /        \  /  \  /    g
  /  \/                  \  /          \/    \/
 /                        \/a           c     e
/                         *81
79
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 78 QA 天眼2007-09-10 16:16:44",
		Text: `
                      4
                      /\                   *8
                     /  \                  /\
0                   /    \         6      /  \
\       2          /      \        /\    /    \
 \      /\        /        \      /  \  /      \
  \    /  \      /          \    /    \/        \     10
   \  /    \    /            \  /     7          \    /\
    \/      \  /              \/                  \  /  \
    1        \/               5                    \/    \
             *3                                    9      \
                                                           \
                                                            \ 11
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 79 Fig 1-2, //Fuzzy",
		Text: `
\0                        4
 \                        /\                   8
  \                      /  \                  /\      10
   \                    /    \         6      /  \    /\
    \       2          /      \        /\    /    \  /  \
     \      /\        /        \      /  \  /      \/    \
      \    /  \      /          \    /    \/       9      \
       \  /    \    /            \  /     7                \
        \/      \  /              \/                        \
        1        \/               5                          \
                 *3                                           \11
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 79 Fig 1-2",
		Text: `
\0                        4
 \                        /\                   *8
  \                      /  \                  /\          10
   \                    /    \         6      /  \        /\
    \       2          /      \        /\    /    \      /  \
     \      /\        /        \      /  \  /      \    /    \
      \    /  \      /          \    /    \/        \  /      \
       \  /    \    /            \  /     7          \/        \
        \/      \  /              \/                  9         \
        1        \/               5                              \
                 *3                                               \11
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 79 Fig 1-2",
		Text: `
\0                        4
 \                        /\                   8
  \                      /  \                  /\     *10
   \                    /    \         6      /  \    /\
    \       2          /      \        /\    /    \  /  \
     \      /\        /        \      /  \  /      \/    \
      \    /  \      /          \    /    \/       9      \
       \  /    \    /            \  /     7                \    12
        \/      \  /              \/                        \  /\
        1        \/               5                          \/  \
                 *3                                           11  \13
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 79 Fig 1-2, //Fuzzy Fuzzy",
		Text: `
\0                        4
 \                        /\                   8
  \                      /  \                  /\
   \                    /    \         6      /  \
    \       2          /      \        /\    /    \  10
     \      /\        /        \      /  \  /      \/\
      \    /  \      /          \    /    \/       9  \
       \  /    \    /            \  /     7            \
        \/      \  /              \/                    \
        1        \/               5                      \
                  3                                       \11
    `,
	},
	test_tdatas_pair{
		Desc: "Lesson 81 QA 袖手旁观 2007-09-19 16:17:15",
		Text: `
                                                  9/
                                                  /
                              *5                 /
                              /\      7         /
                             /  \    /\        /
        1                   /    \  /  \      /
       /\        3         /      \/    \    /
      /  \      /\        /        6     \  /
     /    \    /  \      /                \/
    /      \  /    \    /                  8
   /        \/      \  /
  /          2       \/
 /                    4
/0
    `,
	},
	test_tdatas_pair{
		Desc: "Case 1, a",
		Text: `
        *
        /\
       /  \  /\
  /\  /    \/  \
 /  \/          \
/
    `,
	},
	test_tdatas_pair{
		Desc: "Case 1, b",
		Text: `
        *
        /\
       /  \
  /\  /    \
 /  \/      \  /\
/            \/  \
                  \
    `,
	},
	test_tdatas_pair{
		Desc: "Case 1, c",
		Text: `
                           /
          /\              /
         /  \            /
    /\  /    \          /
   /  \/      \    /\  /
  /            \  /  \/
 /              \/
/
    `,
	},
	test_tdatas_pair{
		Desc: "Case 1, d",
		Text: `
           *
           /\
          /  \
     /\  /    \          /\
    /  \/      \    /\  /  \
   /            \  /  \/    \
  /              \/          \
 /                            \
/
    `,
	},
	test_tdatas_pair{
		Desc: "Case 1, e // Fuzzy",
		Text: `
          *          /\ *
          /\        /  \/\
         /  \      /      \/\
    /\  /    \/\  /          \
   /  \/        \/
  /
 /
/
    `,
	},
	test_tdatas_pair{
		Desc: "Case 1, f",
		Text: `
                            /
          *          /\    /
          /\        /  \  /
         /  \      /    \/
    /\  /    \/\  /
   /  \/        \/
  /             *
 /
/
    `,
	},
	test_tdatas_pair{
		Desc: "Case 2, a",
		Text: `
          *
          /\
         /  \  /\      /
        /    \/  \  /\/
       /          \/
  /\  /           *
 /  \/
/
    `,
	},
	test_tdatas_pair{
		Desc: "Case 2, b",
		Text: `
          *
          /\
         /  \/\        /
        /      \      /
       /        \  /\/
  /\  /          \/
 /  \/
/
    `,
	},
	test_tdatas_pair{
		Desc: "Case 2, c",
		Text: `
                   /
          /\      /
         /  \  /\/
        /    \/
       /
  /\  /
 /  \/
/
    `,
	},
	test_tdatas_pair{
		Desc: "Case 2, d",
		Text: `
          /\
         /  \  /\/\
        /    \/    \
       /
  /\  /
 /  \/
/
    `,
	},
	test_tdatas_pair{
		Desc: "Case 2, e",
		Text: `
                       /
          /\          /
         /  \      /\/
        /    \/\  /
       /        \/
  /\  /
 /  \/
/
    `,
	},
	test_tdatas_pair{
		Desc: "sh600570, 2015-08-14 10:37 -- 13:49 ",
		Text: `
        *3
        /\    5
   1   /  \  /\    7
  /\  /    \/  \  /\              11
 /  \/      4   \/  \            /\
/    2           6   \      9   /  \
0                     \    /\  /    \                            18
                       \  /  \/      \                          /\  21
                        \/    10      \                        /  \/\
                         8             \                  17  /    20\
                                        \            15  /\  /        \
                                         \      13  /\  /  \/          22
                                          \    /\  /  \/    18
                                           \  /  \/    16
                                            \/    14
                                            *12
    `,
	},
	test_tdatas_pair{
		Desc: "sh600570, 2015-08-17 10:24 -- 08-18 09:35",
		Text: `
                      *
                      /\
                     /  \/\           *                /\
                    /      \          /\          /\  /  \
                   /        \    /\  /  \/\      /  \/    \
                  /          \  /  \/      \    /          \
                 /            \/            \  /
            /\  /             *              \/
       /\  /  \/                             *
  /\  /  \/
 /  \/
/
    `,
	},
	test_tdatas_pair{
		Desc: "sh600570, 2015-11-05 13:02 -- 14:25",
		Text: `
\0
 \      2                   /\8
  \    /\            6     /  \
   \  /  \    4     /\    /    \
    \/    \  /\    /  \  /      \
     1     \/  \  /    \/        \
            3   \/      7         \
                *5                 9
    `,
	},
}

func TestParseSegment(t *testing.T) {
	for i, d := range tests_segments {
		lines, segments := text2Segment([]byte(d.Text))
		td := Tdatas{}
		td.Typing.Line = lines
		td.ParseSegment()
		if !test_line_i_price_type_equal(segments, td.Segment.Data) {
			t.Error(
				"\nExample", i,
				"\nFor", d.Desc,
				"\nText", strings.Replace(d.Text, " ", ".", -1),
				"\nexpected", segments,
				"\ngot", td.Segment.Data,
			)
		}
	}
}
