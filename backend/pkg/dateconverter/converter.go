package dateconverter

import (
	apperrors "montelukast/pkg/error"
	"strconv"
	"strings"
	"time"
)

func ConvertDateToTimeStamp(dateStr string) (string, error) {
	parsedTime, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", apperrors.ErrInvalidDate
	}
	timeStamp := parsedTime.Format("2006-01-02 15:04:05")
	return timeStamp, nil
}

func CompareDate(date1 string, date2 string) error {
	parsedTime1, err := time.Parse("2006-01-02 15:04:05", date1)
	if err != nil {
		return apperrors.ErrInvalidDate
	}
	parsedTime2, err := time.Parse("2006-01-02 15:04:05", date2)
	if err != nil {
		return apperrors.ErrInvalidDate
	}
	if !parsedTime1.Before(parsedTime2) {
		return apperrors.ErrInvalidRangeDate
	}
	return nil
}

func GetCurrentDay() string {
	return strings.ToLower(time.Now().Weekday().String())
}

func GetCurrentTime() string {
	return time.Now().Format(`15:04`)
}


func GetCurrentDate() string {
	currentTime := time.Now()
	rawyy, rawmm, rawdd := currentTime.Date()
	yy := strconv.Itoa(rawyy)
	mm := strconv.Itoa(int(rawmm))
	dd := strconv.Itoa(rawdd)
	
	currentDate := strings.Join([]string{yy, mm, dd}, "-")
	return currentDate
}

