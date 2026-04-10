package docgen

import "time"

// Delay pauses execution for the specified duration.
// Used for crash scenarios where game crashes after N seconds.
func Delay(duration string) {
	d, err := time.ParseDuration(duration)
	if err != nil {
		d = 1 * time.Second // default fallback
	}
	time.Sleep(d)
}