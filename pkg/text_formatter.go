package pkg

import "time"

const (
	nocolor = 0
	red		= 31
	green 	= 32
	yellow	= 33
	blue	= 36
	gray	= 37
)

var (
	baseTimestamp time.Time
	emptyFieldMap FieldMap
)

func init() {
	baseTimestamp = time.Now()
}

func main() {

}