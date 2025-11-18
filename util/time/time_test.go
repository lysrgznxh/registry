package time

import (
	"fmt"
	"testing"
	"time"
)

func TestTransTimeStr2Timestamp(t *testing.T) {
	fmt.Println(TransTimeStr2Timestamp("2019-11-18 09:11:05"))
}

func TestParseStr2TimeStamp(t *testing.T) {
	fmt.Println(ParseStr2TimeStamp("2020-05-08 14:15:25"))
}

func TestTimeFormat(t *testing.T) {
	//年月日时分秒  毫秒  时区
	t.Log(time.Now().Format("2006-01-02T15:04:05.000Z"))

	weekday := int(time.Now().Weekday())
	if weekday == 0 { // 如果是周日（默认0）
		weekday = 7 // 转换为7
	}
	fmt.Printf("今天是星期%d\n", weekday)
}
