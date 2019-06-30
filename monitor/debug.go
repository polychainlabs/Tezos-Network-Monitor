package monitor

import "log"

var debugEnabled = true

func debugf(format string, v ...interface{}) {
	if debugEnabled {
		log.Printf(format, v...)
	}
}
func debugln(v ...interface{}) {
	if debugEnabled {
		log.Println(v...)
	}
}