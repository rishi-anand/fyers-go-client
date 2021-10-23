package utils

import "time"

func ToTime(epoch int64) time.Time {
	return time.Unix(epoch, 0)
}

func ToISTTime(in time.Time) time.Time {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	return in.In(loc)
}

func ToIstTimeFromEpoch(epoch int64) time.Time {
	return ToISTTime(ToTime(epoch))
}
