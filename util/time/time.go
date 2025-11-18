package time

import "time"

var cstZone = time.FixedZone("CST", 8*3600) // 东八
// 2019-11-18 09:11:05
func TransTimeStr2Timestamp(s string) (t int) {
	formatTime, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	if err != nil {
		return 0
	}
	return int(formatTime.Unix())
}

func TimeTimeStr2Timestamp(s string) (t int64) {
	formatTime, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	if err != nil {
		return 0
	}
	return formatTime.Unix()
}

func ParseStr2TimeStamp(ts string) (ts_ int64, err_ error) {
	//2020-05-08 14:15:25
	t, err_ := time.ParseInLocation("2006-01-02 15:04:05", ts, cstZone)
	if err_ != nil {
		return
	}
	ts_ = t.Unix()
	return
}

// 解析区块的字符串时间 2021-05-27 10:01:37.558162337 +0000 UTC 转成时间戳
func ParseBlockTime2TimeStamp(blockTime string) (ts_ int64, err_ error) {
	time1, err := time.Parse("2006-01-02 15:04:05.9999999 +0000 UTC", blockTime)
	if err != nil {
		return 0, err
	}
	return time1.Unix(), nil
}

var YMRDHS_Format = "2006-01-02 15:04:05"
var YMD_Format = "2006-01-02"
var YM_Format = "2006-01"

func CurrentDateString() string {
	now := time.Now()
	return now.Format(YMRDHS_Format)
}

// 时间戳转化为时间
func TimeStampToTimeString(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(YMRDHS_Format)
}

// 时间戳转化为时间
func TimeStampToTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

// 时间转化为字符串
func GetTimeStringFormatTime(t time.Time) string {
	return t.Format(YMRDHS_Format)
}

// 时间转化为字符串
func GetTimeStringFormatTimeDay(t time.Time) string {
	return t.Format(YMD_Format)
}

// 时间转化为字符串
func GetTimeStringFormatTimeMonth(t time.Time) string {
	return t.Format(YM_Format)
}

func FirstDayOfISOWeek(year int, week int) time.Time {
	date := time.Date(year, 0, 0, 0, 0, 0, 0, cstZone)
	isoYear, isoWeek := date.ISOWeek()
	for date.Weekday() != time.Monday { // iterate back to Monday
		date = date.AddDate(0, 0, -1)
		isoYear, isoWeek = date.ISOWeek()
	}
	for isoYear < year { // iterate forward to the first day of the first week
		date = date.AddDate(0, 0, 1)
		isoYear, isoWeek = date.ISOWeek()
	}
	for isoWeek < week { // iterate forward to the first day of the given week
		date = date.AddDate(0, 0, 1)
		isoYear, isoWeek = date.ISOWeek()
	}
	return date
}
