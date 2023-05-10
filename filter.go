package api

import "time"

// FilterCondition defines the criterias that must be met when assuming a
// different record presence during the search.
type FilterCondition struct {
	Path string

	Equals any
	IsNull *bool

	StartsWith *string
	EndsWith   *string
	LikeRegex  *string

	LessThan      *float64
	LessEquals    *float64
	GreaterThan   *float64
	GreaterEquals *float64

	After  *time.Time
	Since  *time.Time
	Before *time.Time
	Until  *time.Time

	Invert        *bool
	CaseSensitive *bool
}
