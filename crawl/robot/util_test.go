package robot

import "testing"

// s xxxxxxbegin^_^end return ^_^ or nil
func TestParseParamBeginEnd(t *testing.T) {
	type date_pair struct {
		src, begin, end, exp string
	}
	tests := []date_pair{
		{"xxxxxxbegin^_^end", "begin", "end", "^_^"},
		{`symbol:"sh600112"`, `symbol:"`, `"`, "sh600112"},
	}
	for i, td := range tests {
		act := string(ParseParamBeginEnd([]byte(td.src), []byte(td.begin), []byte(td.end)))
		if act != td.exp {
			t.Error(
				"For", "case", i,
				"expected", td.exp,
				"got", act,
			)
		}
	}
}

func TestParseParamByte(t *testing.T) {
	type date_pair struct {
		src, name, sep, eq, exp string
	}
	tests := []date_pair{
		{"num:2 total:2 start:160322 16:2", "start", " ", ":", "160322"},
	}
	for i, td := range tests {
		act := string(ParseParamByte([]byte(td.src), []byte(td.name), []byte(td.sep), []byte(td.eq)))
		if act != td.exp {
			t.Error(
				"For", "case", i,
				"expected", td.exp,
				"got", act,
			)
		}
	}
}

func TestParseParamInt(t *testing.T) {
	type date_pair struct {
		src, name, sep, eq string
		exp                int
	}
	tests := []date_pair{
		{"num:2 total:2 start:160322 16:2", "start", " ", ":", 160322},
	}
	for i, td := range tests {
		act := ParseParamInt([]byte(td.src), []byte(td.name), []byte(td.sep), []byte(td.eq), td.exp-1)
		if act != td.exp {
			t.Error(
				"For", "case", i,
				"expected", td.exp,
				"got", act,
			)
		}
	}
}
