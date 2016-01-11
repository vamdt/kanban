package crawl

import "testing"

func TestBuildM30s(t *testing.T) {
	t.Skip("todo m30s")
	//TODO merge test M5s
	// 10:00 [09:00, 10:00)
	// 10:30 [10:00, 10:30)
	// 11:00 [10:30, 11:00)
	// 11:30 [11:00, 11:35)
	// 13:30 [13:00, 13:30)
	// 14:00 [13:30, 14:00)
	// 14:30 [14:00, 14:30)
	// 15:00 [14:30, 15:05)
}

func TestTicks2M1s(t *testing.T) {
	type date_pair struct {
		date []string
		len  int
	}
	tests := []date_pair{
		{[]string{
			"2000-01-01 10:00:00",
			"2000-01-02 10:00:00",
			"2000-01-02 10:01:00",
		}, 3},
	}

	fmt := "2006-01-02 15:04:05"
	for i, l := 0, len(tests); i < l; i++ {
		stock := Stock{}
		stock.Ticks = stringSlice2Ticks(tests[i].date)
		for j := 10; j > 0; j-- {
			stock.Ticks2M1s()
			if len(stock.M1s.Data) != tests[i].len {
				t.Error(
					"For", "case", i, tests[i],
					"expected", "m1s.data len=", tests[i].len,
					"got", len(stock.M1s.Data), stock.M1s.Data,
				)
			}
			for k := 0; k < len(tests[i].date); k++ {
				ts := stock.M1s.Data[k].Time.Format(fmt[:len(tests[i].date[k])])
				if ts != tests[i].date[k] {
					t.Error(
						"For", "case", i, "date[", k, "]",
						"expected", tests[i].date[k],
						"got", ts,
					)
				}
			}
		}
	}
}
