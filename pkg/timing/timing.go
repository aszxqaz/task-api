package timing

import (
	"fmt"
	"time"
)

func Timestamp() int64 {
	return time.Now().UTC().Unix()
}

func Elapsed(t1, t2 int64) string {
	d := int(max(t1, t2) - min(t1, t2))
	hours := int(d / 3600)
	d -= hours * 3600
	mins := int(d / 60)
	secs := d - mins*60
	return fmt.Sprintf("%0*d:%0*d:%0*d", 2, hours, 2, mins, 2, secs)
}

func Format(t int64) string {
	dt := time.Unix(t, 0)
	formatted := fmt.Sprintf("%02d.%02d.%d %02d:%02d:%02d",
		dt.Day(), dt.Month(), dt.Year(),
		dt.Hour(), dt.Minute(), dt.Second(),
	)
	return formatted
}
