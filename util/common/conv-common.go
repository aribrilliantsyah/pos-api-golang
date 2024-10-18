package common

import (
	"database/sql"
	"time"
)

func ConvertNullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func ConvertNullTime(nt sql.NullTime) time.Time {
	if nt.Valid {
		return nt.Time
	}
	return time.Time{}
}

func ConvertNullInt32(ni sql.NullInt32) int32 {
	if ni.Valid {
		return ni.Int32
	}
	return 0
}

func ConvertNullInt64(ni sql.NullInt64) int64 {
	if ni.Valid {
		return ni.Int64
	}
	return 0
}

func ConvertNullFLoat64(ni sql.NullFloat64) float64 {
	if ni.Valid {
		return ni.Float64
	}
	return 0
}
