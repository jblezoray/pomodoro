package tuitea

import "strings"

// Each digit: 5 rows × 4 chars wide (using block characters)
var digitPatterns = [10][5]string{
	// 0
	{"▐██▌", "█  █", "█  █", "█  █", "▐██▌"},
	// 1
	{" ▐█ ", "▐██ ", "  █ ", "  █ ", " ███"},
	// 2
	{"███▌", "   █", " ██▌", "█   ", "████"},
	// 3
	{"███▌", "   █", " ██▌", "   █", "███▌"},
	// 4
	{"█  █", "█  █", "████", "   █", "   █"},
	// 5
	{"████", "█   ", "███▌", "   █", "███▌"},
	// 6
	{"▐██▌", "█   ", "███▌", "█  █", "▐██▌"},
	// 7
	{"████", "   █", "  █ ", " █  ", " █  "},
	// 8
	{"▐██▌", "█  █", "▐██▌", "█  █", "▐██▌"},
	// 9
	{"▐██▌", "█  █", " ███", "   █", "▐██▌"},
}

var colonRows = [5]string{" ", "█", " ", "█", " "}

func renderBigTime(s string) []string {
	rows := make([]string, 5)
	for i, ch := range s {
		var col [5]string
		switch {
		case ch == ':':
			col = colonRows
		case ch >= '0' && ch <= '9':
			col = digitPatterns[ch-'0']
		default:
			for j := range col {
				col[j] = "    "
			}
		}
		sep := " "
		if i == 0 {
			sep = ""
		}
		for r := 0; r < 5; r++ {
			rows[r] += sep + col[r]
		}
	}
	for i := range rows {
		rows[i] = strings.TrimRight(rows[i], " ")
	}
	return rows
}
