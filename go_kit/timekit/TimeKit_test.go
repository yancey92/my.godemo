package timekit

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeDemo(t *testing.T) {
	now := time.Now()
	dateStr, err := TimeToString(now, DateFormat_YYYY_MM_DD)
	if err != nil {
		fmt.Printf("%v\n", err)
		t.Fail()
	}

	fmt.Printf("当前日期:%s\n", dateStr)
}

func TestStringToTime(t *testing.T) {
	tm, err := StringToTime("2016-11-11", DateFormat_YYYY_MM_DD)
	if err != nil {
		fmt.Printf("%v", err)
		t.Fail()
	}
	tmStr, err := TimeToString(tm, DateFormat_YYYY_MM_DD_HH_MM_SS)
	if err != nil {
		fmt.Printf("%v", err)
		t.Fail()
	} else {
		fmt.Println(tmStr) //2016-11-11 00:00:00
	}
}

func TestGetTimeSsAndDate(t *testing.T) {
	ss, dateStr, _ := GetTimeSsAndDate(time.Now(), DateFormat_YYYY_MM_DD_HH_MM_SS)
	fmt.Printf("秒=%d,日期=%s\n", ss, dateStr)
}

func TestGetTimeMsAndDate(t *testing.T) {
	ms, dateStr, _ := GetTimeMsAndDate(time.Now(), DateFormat_YYYY_MM_DD_HH_MM_SS)
	fmt.Printf("毫秒=%d,日期=%s\n", ms, dateStr)
}

func TestGetAfterDayMs(t *testing.T) {
	fmt.Println(GetAfterDayMsAndDate("2016-11-11"))
}

func TestGetEndDayMs(t *testing.T) {
	fmt.Println(GetEndDayMsAndDate("2016-11-11"))
}

//获取某天凌晨时间
func TestGetDayTime(t *testing.T) {
	tm, err := StringToTime("2016-11-11", DateFormat_YYYY_MM_DD)
	if err != nil {
		fmt.Printf("%v", err)
	}
	ms, dateStr, _ := GetTimeMsAndDate(tm, DateFormat_YYYY_MM_DD)
	fmt.Printf("毫秒=%d,日期=%s\n", ms, dateStr)
}

func TestDateRangeVaild(t *testing.T) {
	startDate := "2016-11-11"
	endDate := "2016-11-12"
	result, err := DateRangeVaild(startDate, endDate, DateFormat_YYYY_MM_DD)
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("开始日期=%s 结束日期=%s 范围是否合法 %v\n", startDate, endDate, result)
}

func TestDateRangeVaild2(t *testing.T) {
	startDate := "2016-11-11 11:11:11"
	endDate := "2016-11-11 23:00:00"
	result, err := DateRangeVaild(startDate, endDate, DateFormat_YYYY_MM_DD_HH_MM_SS)
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("开始日期=%s 结束日期=%s 范围是否合法 %v\n", startDate, endDate, result)
}

func TestDefaultDateRangeVaild(t *testing.T) {
	startDate := "2016-11-11"
	endDate := "2016-11-11"
	result, err := DefaultDateRangeVaild(startDate, endDate)
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("开始日期=%s 结束日期=%s 范围是否合法 %v\n", startDate, endDate, result)
}

func TestCheckDateRange(t *testing.T) {
	checkDate := "2016-11-12 13:11:10"
	startDate := "2016-11-11 11:11:11"
	endDate := "2016-11-11 23:00:00"
	result, err := CheckDateRange(checkDate, startDate, endDate, DateFormat_YYYY_MM_DD_HH_MM_SS)
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("校验日期=%s 开始日期=%s 结束日期=%s 范围是否合法 %v\n", checkDate, startDate, endDate, result)
}

func TestDefaultCheckDateRange(t *testing.T) {
	checkDate := "2016-11-11"
	startDate := "2016-11-11"
	endDate := "2016-11-11"
	result, err := DefaultCheckDateRange(checkDate, startDate, endDate)
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("校验日期=%s 开始日期=%s 结束日期=%s 范围是否合法 %v\n", checkDate, startDate, endDate, result)
}

func TestGetDateByTime(t *testing.T) {
	tt := "2016-11-11 11:12:13"
	result, err := GetDateByTime(tt, DateFormat_YYYY_MM_DD_HH_MM_SS)
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("格式化后日期 %s\n", result)
}

func TestMs2Date(t *testing.T) {
	ms := 0
	sec := int64(ms / 1000)
	result, err := TimeToString(time.Unix(sec, 0), DateFormat_YYYY_MM_DD_HH_MM_SS)
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("格式化后日期 %s\n", result)
}

func TestDateSubMonth(t *testing.T) {
	startDate := "2017-11-11"
	endDate := "2018-01-11"
	result, err := DateSubMonth(startDate, endDate)
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("开始日期=%s 结束日期=%s 结果= %v\n", startDate, endDate, result)
}
func TestDateStrSplit(t *testing.T) {
	s := "2017-12-13 12:28:55"
	date, time, err := DateStrSplit(s)
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Println(s)
	fmt.Printf("日期：%s  时间：%s\n", date, time)
}

func TestDateStrConv(t *testing.T) {
	s := "2017-12-13 12:28:55"
	r, err := DateStrConv(s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r)
}
