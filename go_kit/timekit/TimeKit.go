// @Title 处理日常使用的时间与格式化日期字符串之间的转换工具

package timekit

import (
	"errors"
	"fmt"
	"gitlab.gumpcome.com/common/go_kit/logiccode"
	"gitlab.gumpcome.com/common/go_kit/strkit"
	"strconv"
	"time"
)

const (
	DateFormat_YYYY_MM_DD_HH_MM_SS_MS string = "2006-01-02 15:04:05 000"
	DateFormat_YYYY_MM_DD_HH_MM_SS    string = "2006-01-02 15:04:05"
	DateFormat_YYYY_MM_DD             string = "2006-01-02"
	DateFormat_YYYYMMDDHHMMSSMS       string = "20060102150405000"
	DateFormat_YYYYMMDDHHMMSS         string = "20060102150405"
	DateFormat_YYYYMMDD               string = "20060102"
	DateFormat_YYYYMM                 string = "200601"
	DATE_RANGE_BEFORE                 string = "before"
	DATE_RANGE_MIDDLE                 string = "middle"
	DATE_RANGE_AFTER                  string = "after"
	DATE_RANGE_ERROR                  string = "error"
)

// @Title 格式化时间对象为日期字符串格式
// @Description
// @param dateTime  time.Time
// @param dateStyle string
// usage:
//	TimeToString(time.Now(), DateFormat_YYYY_MM_DD_HH_MM_SS)
func TimeToString(time time.Time, dateStyle string) (string, error) {
	if dateStyle == "" {
		return "", errors.New("dateStyle is empty")
	}
	switch dateStyle {
	case DateFormat_YYYY_MM_DD_HH_MM_SS,
		DateFormat_YYYY_MM_DD,
		DateFormat_YYYYMMDDHHMMSS,
		DateFormat_YYYYMMDD,
		DateFormat_YYYYMM:
		return time.Format(dateStyle), nil
	case DateFormat_YYYY_MM_DD_HH_MM_SS_MS:
		return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d %03d", time.Year(), time.Month(), time.Day(), time.Hour(), time.Minute(), time.Second(), time.Nanosecond()/1000000), nil
	case DateFormat_YYYYMMDDHHMMSSMS:
		return fmt.Sprintf("%04d%02d%02d%02d%02d%02d%03d", time.Year(), time.Month(), time.Day(), time.Hour(), time.Minute(), time.Second(), time.Nanosecond()/1000000), nil
	default:
		return "", errors.New("dateStyle is error")
	}
}

// @Title 获取某时刻的毫秒值和日期
// @Description
// @param dateTime  time.Time
// @param dateStyle string
// usage:
//	GetTimeMsAndDate(time.Now(), DateFormat_YYYY_MM_DD_HH_MM_SS)
func GetTimeMsAndDate(dateTime time.Time, dateStyle string) (int64, string, error) {
	ms := dateTime.UnixNano() / 1000000
	date, err := TimeToString(dateTime, dateStyle)
	return ms, date, err
}

// @Title 获取某时刻的秒值和日期
// @Description
// @param dateTime  time.Time
// @param dateStyle string
// usage:
//	GetTimeSsAndDate(time.Now(), DateFormat_YYYY_MM_DD_HH_MM_SS)
func GetTimeSsAndDate(dateTime time.Time, dateStyle string) (int64, string, error) {
	ss := dateTime.Unix()
	date, err := TimeToString(dateTime, dateStyle)
	return ss, date, err
}

// @Title 获取当前时刻的毫秒值和日期
// @Description
// @param dateStyle string
// usage:
//	GetNowTimeMsAndDate(DateFormat_YYYY_MM_DD_HH_MM_SS)
func GetNowTimeMsAndDate(dateStyle string) (int64, string, error) {
	return GetTimeMsAndDate(time.Now(), dateStyle)
}

// @Title 获取当前时刻的秒值和日期
// @Description
// @param dateStyle string
// usage:
//	GetNowTimeSsAndDate(DateFormat_YYYY_MM_DD_HH_MM_SS)
func GetNowTimeSsAndDate(dateStyle string) (int64, string, error) {
	return GetTimeSsAndDate(time.Now(), dateStyle)
}

// @Title 日期字符串格式转换成时间对象
// @Description
// @param dateTime  string
// @param dateStyle string
// usage:
//	注意:不使用毫秒类日期格式
//	StringToTime("2016-11-11", DateFormat_YYYY_MM_DD)
func StringToTime(date string, dateStyle string) (time.Time, error) {
	if date == "" {
		return time.Time{}, errors.New("date is empty")
	}

	if dateStyle == "" {
		return time.Time{}, errors.New("dateStyle is empty")
	}

	switch dateStyle {
	case DateFormat_YYYY_MM_DD_HH_MM_SS,
		DateFormat_YYYY_MM_DD,
		DateFormat_YYYYMMDDHHMMSS,
		DateFormat_YYYYMMDD,
		DateFormat_YYYYMM:
		return time.ParseInLocation(dateStyle, date, time.Local)
	default:
		return time.Time{}, errors.New("dateStyle is error")
	}
}

// @Title 获取某日的后一天零时该时刻的毫秒值和日期
// @Description 例如:date=2016-11-11 返回 1478880000000 2016-11-12 00:00:00
// @param date string
// usage:
//	date 必须是YYYY-MM-DD格式
//	GetAfterDayMsAndDate("2016-11-11")
func GetAfterDayMsAndDate(date string) (int64, string, error) {
	dateTime, err := StringToTime(date, DateFormat_YYYY_MM_DD)
	if err != nil {
		return 0, "", err
	}
	afterDateTime := dateTime.Add(time.Second * 60 * 60 * 24)
	return GetTimeMsAndDate(afterDateTime, DateFormat_YYYY_MM_DD_HH_MM_SS)
}

// @Title 获取某日最后一分钟该时刻的毫秒值和日期
// @Description 例如:date=2016-11-11 返回 1478879999999 2016-11-11 23:59:59
// @param date string
// usage:
//	date 必须是YYYY-MM-DD格式
//	GetEndDayMsAndDate("2016-11-11")
func GetEndDayMsAndDate(date string) (int64, string, error) {
	dateTime, err := StringToTime(date, DateFormat_YYYY_MM_DD)
	if err != nil {
		return 0, "", err
	}
	afterDateTime := dateTime.Add(time.Second*60*60*24 - 1)
	return GetTimeMsAndDate(afterDateTime, DateFormat_YYYY_MM_DD_HH_MM_SS)
}

// @Title 根据YYYY-MM-DD HH:MM:SS格式的时间获取对应的YYYY-MM-DD格式日期
// @Description 例如:date=2016-11-11 12:00:00 返回 2016-11-11
// @param timeStr string
// usage:
//	GetDateByTime("2016-11-11 12:00:00")
func GetDateByTime(timeStr string, dateStyle string) (string, error) {
	t, err := StringToTime(timeStr, dateStyle)
	if err != nil {
		return "", err
	}
	_, d, err := GetTimeMsAndDate(t, DateFormat_YYYY_MM_DD)
	if err != nil {
		return "", err
	}
	return d, nil
}

// @Title 判断checkDate是否在startDate和endDate日期范围之间
// usage:
//	DATE_RANGE_BEFORE: checkDate在startDate前面
//	DATE_RANGE_AFTER: checkDate在endDate前面
//	DATE_RANGE_MIDDLE: checkDate在startDate和endDate日期范围之间
func CheckDateRange(checkDate string, startDate string, endDate string, dateStyle string) (string, error) {
	endDateTime, err := StringToTime(endDate, dateStyle)
	if err != nil {
		return DATE_RANGE_ERROR, err
	}
	startDateTime, err := StringToTime(startDate, dateStyle)
	if err != nil {
		return DATE_RANGE_ERROR, err
	}
	checkDateTime, err := StringToTime(checkDate, dateStyle)
	if err != nil {
		return DATE_RANGE_ERROR, err
	}
	if !startDateTime.Before(endDateTime) {
		return DATE_RANGE_ERROR, logiccode.New(100302, "data range error")
	}
	if checkDateTime.Before(startDateTime) {
		return DATE_RANGE_BEFORE, nil
	}
	if checkDateTime.After(endDateTime) {
		return DATE_RANGE_AFTER, nil
	}
	return DATE_RANGE_MIDDLE, nil
}

func DefaultCheckDateRange(checkDate string, startDate string, endDate string) (string, error) {
	endDateTime, err := StringToTime(endDate, DateFormat_YYYY_MM_DD)
	if err != nil {
		return DATE_RANGE_ERROR, err
	}
	afterEndDateTime := endDateTime.Add(time.Second*60*60*24 - 1)
	startDateTime, err := StringToTime(startDate, DateFormat_YYYY_MM_DD)
	if err != nil {
		return DATE_RANGE_ERROR, err
	}
	checkDateTime, err := StringToTime(checkDate, DateFormat_YYYY_MM_DD)
	if err != nil {
		return DATE_RANGE_ERROR, err
	}
	if !startDateTime.Before(afterEndDateTime) {
		return DATE_RANGE_ERROR, logiccode.New(100302, "data range error")
	}
	if checkDateTime.Before(startDateTime) {
		return DATE_RANGE_BEFORE, nil
	}
	if checkDateTime.After(afterEndDateTime) {
		return DATE_RANGE_AFTER, nil
	}
	return DATE_RANGE_MIDDLE, nil
}

// @Title 校验startDate和endDate日期范围的有效性
func DateRangeVaild(startDate string, endDate string, dateStyle string) (bool, error) {
	endDateTime, err := StringToTime(endDate, dateStyle)
	if err != nil {
		return false, err
	}
	startDateTime, err := StringToTime(startDate, dateStyle)
	if err != nil {
		return false, err
	}
	if startDateTime.Before(endDateTime) {
		return true, nil
	}
	return false, err
}

func DefaultDateRangeVaild(startDate string, endDate string) (bool, error) {
	endDateTime, err := StringToTime(endDate, DateFormat_YYYY_MM_DD)
	if err != nil {
		return false, err
	}
	afterEndDateTime := endDateTime.Add(time.Second*60*60*24 - 1)
	startDateTime, err := StringToTime(startDate, DateFormat_YYYY_MM_DD)
	if err != nil {
		return false, err
	}
	if startDateTime.Before(afterEndDateTime) {
		return true, nil
	}
	return false, err
}

// @Title 获取两个日期之内的时间差，返回YYYYMM类型的字符串数组
// startDate 开始日期 YYYY-MM-DD 格式
// endDate 结束日期 YYYY-MM-DD 格式
func DateSubMonth(startDate, endDate string) ([]string, error) {
	start, _ := StringToTime(startDate, DateFormat_YYYY_MM_DD)
	end, _ := StringToTime(endDate, DateFormat_YYYY_MM_DD)
	dateSlice := make([]string, 0)
	if start.Year() < end.Year() {
		for i := start.Year(); i <= end.Year(); i++ {
			if i == start.Year() {
				for j := start.Month(); j <= 12; j++ {
					t := time.Date(i, j, 1, 0, 0, 0, 0, time.Local)
					timeString, err := TimeToString(t, DateFormat_YYYYMM)
					if err != nil {
						return nil, err
					}
					dateSlice = append(dateSlice, timeString)
				}
			}
			if i > start.Year() && i < end.Year() {
				for j := time.January; j <= 12; j++ {
					t := time.Date(i, j, 1, 0, 0, 0, 0, time.Local)
					timeString, err := TimeToString(t, DateFormat_YYYYMM)
					if err != nil {
						return nil, err
					}
					dateSlice = append(dateSlice, timeString)
				}
			}
			if i == end.Year() {
				for j := time.January; j <= end.Month(); j++ {
					t := time.Date(i, j, 1, 0, 0, 0, 0, time.Local)
					timeString, err := TimeToString(t, DateFormat_YYYYMM)
					if err != nil {
						return nil, err
					}
					dateSlice = append(dateSlice, timeString)
				}
			}
		}
	}
	if start.Year() == end.Year() {
		for j := start.Month(); j <= end.Month(); j++ {
			t := time.Date(start.Year(), j, 1, 0, 0, 0, 0, time.Local)
			timeString, err := TimeToString(t, DateFormat_YYYYMM)
			if err != nil {
				return nil, err
			}
			dateSlice = append(dateSlice, timeString)
		}
	}
	return dateSlice, nil
}

func JoinDateTrix(startTime, endTime string) ([]string, error) {
	s := []string{startTime, endTime}
	t := make([]time.Time, 0, 2)
	for _, i := range s {
		r, err := time.Parse(DateFormat_YYYY_MM_DD, i)
		if err != nil {
			rr, err2 := time.Parse(DateFormat_YYYY_MM_DD_HH_MM_SS, i)
			if err2 != nil {
				return nil, errors.New(fmt.Sprintf("日期时间参数 %s 格式错误 ，合法格式为：2006-01-02 15:04:05 或2006-01-02", i))
			}
			t = append(t, rr)
		} else {
			t = append(t, r)
		}
	}
	sn := t[0].Year()*100 + int(t[0].Month())
	en := t[1].Year()*100 + int(t[1].Month())
	if sn > en {
		return nil, errors.New(fmt.Sprintf("参数不合法：开始时间 %s 晚于结束时间 %s", startTime, endTime))
	}
	mm := int(t[0].Month())
	yy := t[0].Year() * 100
	result := make([]string, 0)
	for i := sn; i <= en; {
		result = append(result, strconv.Itoa(i))
		mm++
		if mm > 12 {
			mm = 1
			yy += 100
		}
		i = yy + mm
	}
	return result, nil
}

//解析 "2006-01-02 15:04:05" 格式时间字符串为两部分 "20060102" "150405"
func DateStrSplit(date string) (string, string, error) {
	dt, err := StringToTime(date, DateFormat_YYYY_MM_DD_HH_MM_SS)
	if err != nil {
		return "", "", err
	}
	s, err := TimeToString(dt, DateFormat_YYYYMMDDHHMMSS)
	if err != nil {
		return "", "", err
	}
	return strkit.SubStr(s, 0, 8), strkit.SubStr(s, 8, 14), nil
}

// 时间字符串转换，"2006-01-02 15:04:05" 转换为 "20060102150405"
func DateStrConv(date string) (string, error) {
	t, err := time.Parse(DateFormat_YYYY_MM_DD_HH_MM_SS, date)
	if err != nil {
		return "", err
	}
	return t.Format(DateFormat_YYYYMMDDHHMMSS), nil
}
