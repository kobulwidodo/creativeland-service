package timeutils

import (
	"fmt"
	"time"
)

func DiffForHumans(t time.Time) string {
	now := time.Now()
	duration := now.Sub(t)

	switch {
	case duration.Hours() < 1:
		return fmt.Sprintf("%d menit yang lalu", int(duration.Minutes()))
	case duration.Hours() < 24:
		return fmt.Sprintf("%d jam yang lalu", int(duration.Hours()))
	case duration.Hours() < 24*30:
		return fmt.Sprintf("%d hari yang lalu", int(duration.Hours()/24))
	case duration.Hours() < 24*30*12:
		return fmt.Sprintf("%d bulan yang lalu", int(duration.Hours()/(24*30)))
	default:
		return fmt.Sprintf("%d tahun yang lalu", int(duration.Hours()/(24*30*12)))
	}
}
