package api

import "fmt"

type Message struct {
	Room int
	Type string
	Data string
}

func FormatFloat(f float64) string {
	return fmt.Sprintf("%.2f", f)
}
