package gw_date

import (
	"time"
)

func GetLastDayOfMonth(targetTime time.Time, location *time.Location) time.Time {
	firstOfMonth := time.Date(targetTime.Year(), targetTime.Month(), 1, 0, 0, 0, 0, location)
	firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)
	lastOfMonth := firstOfNextMonth.Add(-time.Duration(1) * time.Second)
	return lastOfMonth
}
func GetFirstDayOfMonth(targetTime time.Time, location *time.Location) time.Time {
	firstOfMonth := time.Date(targetTime.Year(), targetTime.Month(), 1, 0, 0, 0, 0, location)
	return firstOfMonth
}
func GetLastDayOfNextMonth(targetTime time.Time, location *time.Location) time.Time {
	firstOfMonth := time.Date(targetTime.Year(), targetTime.Month(), 1, 0, 0, 0, 0, location)
	firstOf2MonthAfter := firstOfMonth.AddDate(0, 2, 0)
	lastOfNextMonth := firstOf2MonthAfter.Add(-time.Duration(1) * time.Second)
	return lastOfNextMonth
}
func GetFirstDayOfNextMonth(targetTime time.Time, location *time.Location) time.Time {
	firstOfMonth := time.Date(targetTime.Year(), targetTime.Month(), 1, 0, 0, 0, 0, location)
	firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)
	return firstOfNextMonth
}

//当月末を取得
func TruncDateAndGetLastSecondOfMonth(targetTime time.Time) (time.Time, time.Time) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	truncatedDate := targetTime.In(jst)
	truncatedDate = truncatedDate.Truncate(time.Hour).Add(-time.Duration(truncatedDate.Hour()) * time.Hour)
	//今月末を取得
	toDate := GetLastDayOfMonth(truncatedDate, jst)
	return truncatedDate, toDate
}

//月末の更新日（月末日 AM 9:00）を考慮して、最終日当日(AM0:00 以降)の場合は、翌月末を更新日に設定する
// return truncateDate, toDate, trialEndDate
func TruncDateAndGetTrialEndAndToDate(targetTime time.Time) (time.Time, time.Time, *time.Time) {
	truncatedDate, toDate := TruncDateAndGetLastSecondOfMonth(targetTime)
	if truncatedDate.Format("20060102") == toDate.Format("20060102") {
		//toDateを翌月末まで伸ばし、trial endはその月の月初とする
		jst, _ := time.LoadLocation("Asia/Tokyo")
		trialEnd := GetFirstDayOfNextMonth(toDate, jst)
		return truncatedDate,
			GetLastDayOfNextMonth(toDate, jst),
			&trialEnd
	} else {
		return truncatedDate, toDate, nil
	}
}

//X日後の日本時間ラスト秒を取得
func GetLastSecondOfDay(targetTime *time.Time, addDays int) *time.Time {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	truncatedDate := targetTime.In(jst)
	truncatedDate = truncatedDate.Truncate(time.Hour).Add(-time.Duration(truncatedDate.Hour()) * time.Hour)
	//今月末を取得
	truncatedDate = truncatedDate.AddDate(0, 0, addDays+1)
	retDate := truncatedDate.Add(-time.Duration(1) * time.Second)
	return &retDate
}
