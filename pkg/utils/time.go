package utils

import (
	"doan/pkg/constants"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"time"
)

func GetCurrentTimeUnix() int64 {
	return time.Now().Unix()
}

func GetCurrentTime() time.Time {
	return time.Now()
}

func GetCurrentTimeString() string {
	return time.Now().Format(time.RFC3339)
}

func GetCurrentTimeStringFormat(format string) string {
	return time.Now().Format(format)
}

func ConvertTimeToString(t time.Time) string {
	return t.Format(time.RFC3339)
}

func GetCurrentTimeVietNam() (time.Time, error) {
	// Lấy múi giờ Việt Nam
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		return time.Time{}, err
	}
	// Lấy thời gian hiện tại theo múi giờ Việt Nam
	currentTime := time.Now().In(loc)
	return currentTime, nil
}

func GetCurrentTimeVietNamString() (string, error) {
	// Lấy thời gian hiện tại theo múi giờ Việt Nam
	currentTime, err := GetCurrentTimeVietNam()
	if err != nil {
		return "", err
	}
	return currentTime.Format(time.RFC3339), nil
}

func ConvertStringToDateOnly(dateStr string) (time.Time, error) {
	result, err := time.Parse(string(constants.DateOnly), dateStr)
	return result, err
}

func ConvertTimeToVietNamDateString(date time.Time) string {
	return date.Format(string(constants.DateOnly))
}

func GetTimeStamppb(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func GetTimeStamppbFromUnix(timeUnix int64) *timestamppb.Timestamp {
	dateTime := time.Unix(timeUnix, 0)
	return timestamppb.New(dateTime)
}
func ToGormDeletedAt(t *time.Time) *gorm.DeletedAt {
	if t == nil {
		return nil
	}
	return &gorm.DeletedAt{
		Time:  *t,
		Valid: true,
	}
}

func IsValidTime(value string, format string) bool {
	if _, err := time.Parse(format, value); err == nil {
		return true
	}
	return false
}
