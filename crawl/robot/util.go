package robot

import (
	"bytes"
	"strconv"
)

// s xxxxxxbegin^_^end return ^_^ or nil
func ParseParamBeginEnd(s, begin, end []byte) []byte {
	i := bytes.Index(s, begin)
	if i < 0 {
		return nil
	}
	s = s[i+len(begin):]

	if end == nil {
		return s
	}

	i = bytes.Index(s, end)
	if i < 0 {
		return nil
	}
	return s[:i]
}

// s a=b&c=d&e=f name=c sep=& eq== return d or nil
func ParseParamByte(s, name, sep, eq []byte) []byte {
	lines := bytes.Split(s, sep)
	for i, _ := range lines {
		if !bytes.HasPrefix(lines[i], name) {
			continue
		}
		v := bytes.Split(lines[i], eq)
		if len(v) < 2 {
			continue
		}
		if bytes.Equal(v[0], name) {
			return v[1]
		}
	}
	return nil
}

// like ParseParamByte return int
func ParseParamInt(s, name, sep, eq []byte, defv int) int {
	b := ParseParamByte(s, name, sep, eq)
	if len(b) > 0 {
		i, _ := strconv.Atoi(string(b))
		return i
	}
	return defv
}
