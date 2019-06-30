package alert

import "time"

// Throttle Alerts
var alertThrottle map[string]*time.Time

// return true if an alert has already been sent with this message in the
// provided duration
func hasAlreadyAlerted(message string, duration time.Duration) bool {
	if alertThrottle == nil {
		alertThrottle = map[string]*time.Time{}
	}

	lastAlert, _ := alertThrottle[message]
	now := time.Now()
	if lastAlert == nil || now.Sub(*lastAlert) > duration {
		alertThrottle[message] = &now
		return false
	}
	return true
}
