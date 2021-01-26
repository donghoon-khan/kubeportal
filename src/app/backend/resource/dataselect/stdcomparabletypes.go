package dataselect

import (
	"strings"
	"time"
)

type StdComparableInt int

func (stdComparableInt StdComparableInt) Compare(otherV ComparableValue) int {
	other := otherV.(StdComparableInt)
	return intsCompare(int(stdComparableInt), int(other))
}

func (stdComparableInt StdComparableInt) Contains(otherV ComparableValue) bool {
	return stdComparableInt.Compare(otherV) == 0
}

type StdComparableString string

func (stdComparableString StdComparableString) Compare(otherV ComparableValue) int {
	other := otherV.(StdComparableString)
	return strings.Compare(string(stdComparableString), string(other))
}

func (stdComparableString StdComparableString) Contains(otherV ComparableValue) bool {
	other := otherV.(StdComparableString)
	return strings.Contains(string(stdComparableString), string(other))
}

type StdComparableRFC3339Timestamp string

func (stdComparableRFC3339Timestamp StdComparableRFC3339Timestamp) Compare(otherV ComparableValue) int {
	other := otherV.(StdComparableRFC3339Timestamp)
	selfTime, err1 := time.Parse(time.RFC3339, string(stdComparableRFC3339Timestamp))
	otherTime, err2 := time.Parse(time.RFC3339, string(other))

	if err1 != nil || err2 != nil {
		return strings.Compare(string(stdComparableRFC3339Timestamp), string(other))
	}
	return ints64Compare(selfTime.Unix(), otherTime.Unix())
}

func (stdComparableRFC3339Timestamp StdComparableRFC3339Timestamp) Contains(otherV ComparableValue) bool {
	return stdComparableRFC3339Timestamp.Compare(otherV) == 0
}

type StdComparableTime time.Time

func (stdComparableTime StdComparableTime) Compare(otherV ComparableValue) int {
	other := otherV.(StdComparableTime)
	return ints64Compare(time.Time(stdComparableTime).Unix(), time.Time(other).Unix())
}

func (stdComparableTime StdComparableTime) Contains(otherV ComparableValue) bool {
	return stdComparableTime.Compare(otherV) == 0
}

func intsCompare(a, b int) int {
	if a > b {
		return 1
	} else if a == b {
		return 0
	}
	return -1
}

func ints64Compare(a, b int64) int {
	if a > b {
		return 1
	} else if a == b {
		return 0
	}
	return -1
}
