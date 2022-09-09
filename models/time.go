/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	log "github.com/sirupsen/logrus"
	"time"
)

const SvgTimeFormat = "15:04:05 02-01-2006 (UTC)"

//MsToSeconds millisecond to second
func MsToSeconds(millisecond int64) int64 {
	return time.UnixMilli(millisecond).Unix()
}

func NanoToSeconds(nano int64) int64 {
	if nano == 0 {
		return 0
	}
	return nano / 1e9
}

func NanoToMs(nano int64) int64 {
	if nano == 0 {
		return 0
	}
	return nano / 1e6
}

func GetDateDiffFromNow(layout string, findTime string, offset int) int64 {
	var diff int64
	tz := time.Unix(GetNowTimeUnix(), 0)
	today := time.Date(tz.Year(), tz.Month(), tz.Day(), 0, 0, 0, 0, tz.Location())

	if layout != "2006-01-02" && layout != "2006-01" && layout != "2006" {
		return 0
	}
	t2, err := time.ParseInLocation(layout, findTime, time.Local)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Get Date Part From Now ParseInLocation Failed")
		return 0
	}
	switch layout {
	case "2006-01-02":
		t1 := today.AddDate(0, 0, offset)
		diff = int64(t1.Sub(t2).Hours() / 24)
		if diff > 0 {
			if int(t1.Sub(t2).Milliseconds())%86400000 > int(86400000-t2.Unix()%86400000) {
				diff += 1
			}
			return diff
		}

	case "2006-01":
		t1 := today.AddDate(0, offset, 0)
		d1 := t1.Day()
		d2 := t2.Day()
		m1 := int64(t1.Month())
		m2 := int64(t2.Month())
		yearInterval := m1 - m2
		if m1 < m2 || m1 == m2 && d1 < d2 {
			yearInterval--
		}
		monthInterval := (m1 + 12) - m2
		if d1 < d2 {
			monthInterval--
		}
		monthInterval %= 12
		month := yearInterval*12 + monthInterval
		if month < 0 {
			return 0
		}
		return month

	case "2006":
		t1 := today.AddDate(0, offset, 0)
		d1 := t1.Day()
		d2 := t2.Day()
		m1 := int64(t1.Month())
		m2 := int64(t2.Month())
		yearInterval := m1 - m2
		if m1 < m2 || m1 == m2 && d1 < d2 {
			yearInterval--
		}
		monthInterval := (m1 + 12) - m2
		if d1 < d2 {
			monthInterval--
		}
		monthInterval %= 12
		year := (yearInterval*12 + monthInterval) / 12
		if year < 0 {
			return 0
		}
		return year
	}

	return diff
}
