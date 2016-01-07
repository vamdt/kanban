package crawl

import (
	"strconv"
	"testing"
	"time"
)

func TestAtoi(t *testing.T) {
	i, _ := strconv.Atoi("-05")
	if i != -5 {
		t.Error(
			"For", "-05",
			"expected", -5,
			"got", i,
		)
	}
}

func TestByteString(t *testing.T) {
	var b []byte
	s := string(b)
	if s != "" {
		t.Error(
			"For", "",
			"expected", "",
			"got", s,
		)
	}
}

func TestTimeNowHour(t *testing.T) {
	now := time.Now().UTC()
	t.Logf("now utc hour %d", now.Hour())
}

func TestTimeBefore(t *testing.T) {
	now := time.Now()
	before := now.Add(-1 * time.Second)
	if now.Before(now) {
		t.Error(
			"For", "now Before now",
			"expected", false,
			"got", true,
		)
	}

	if now.Before(now) {
		t.Error(
			"For", "now Before now",
			"expected", true,
			"got", false,
		)
	}

	if now.Before(before) {
		t.Error(
			"For", "now Before before",
			"expected", true,
			"got", false,
		)
	}
}