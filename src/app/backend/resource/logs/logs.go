package logs

import (
	"sort"
	"strings"
)

var LineIndexNotFound = -1

var DefaultDisplayNumLogLines = 100

var MaxLogLines int = 2000000000

const (
	NewestTimestamp = "newest"
	OldestTimestamp = "oldest"
)

const (
	Beginning = "beginning"
	End       = "end"
)

var NewestLogLineId = LogLineId{
	LogTimestamp: NewestTimestamp,
}

var OldestLogLineId = LogLineId{
	LogTimestamp: OldestTimestamp,
}

var DefaultSelection = &Selection{
	OffsetFrom:      1 - DefaultDisplayNumLogLines,
	OffsetTo:        1,
	ReferencePoint:  NewestLogLineId,
	LogFilePosition: End,
}

var AllSelection = &Selection{
	OffsetFrom:     -MaxLogLines,
	OffsetTo:       MaxLogLines,
	ReferencePoint: NewestLogLineId,
}

type LogDetails struct {
	Info      LogInfo `json:"info"`
	Selection `json:"selection"`
	LogLines  `json:"logs"`
}

type LogInfo struct {
	PodName           string       `json:"podName"`
	ContainerName     string       `json:"containerName"`
	InitContainerName string       `json:"initContainerName"`
	FromDate          LogTimestamp `json:"fromDate"`
	ToDate            LogTimestamp `json:"toDate"`
	Truncated         bool         `json:"truncated"`
}

type Selection struct {
	ReferencePoint  LogLineId `json:"referencePoint"`
	OffsetFrom      int       `json:"offsetFrom"`
	OffsetTo        int       `json:"offsetTo"`
	LogFilePosition string    `json:"logFilePosition"`
}

type LogLineId struct {
	LogTimestamp `json:"timestamp"`
	LineNum      int `json:"lineNum"`
}

type LogLines []LogLine

type LogLine struct {
	Timestamp LogTimestamp `json:"timestamp"`
	Content   string       `json:"content"`
}

type LogTimestamp string

func (self LogLines) SelectLogs(logSelection *Selection) (LogLines, LogTimestamp, LogTimestamp, Selection, bool) {
	requestedNumItems := logSelection.OffsetTo - logSelection.OffsetFrom
	referenceLineIndex := self.getLineIndex(&logSelection.ReferencePoint)
	if referenceLineIndex == LineIndexNotFound || requestedNumItems <= 0 || len(self) == 0 {
		return LogLines{}, "", "", Selection{}, false
	}
	fromIndex := referenceLineIndex + logSelection.OffsetFrom
	toIndex := referenceLineIndex + logSelection.OffsetTo
	lastPage := false
	if requestedNumItems > len(self) {
		fromIndex = 0
		toIndex = len(self)
		lastPage = true
	} else if toIndex > len(self) {
		fromIndex -= toIndex - len(self)
		toIndex = len(self)
		lastPage = logSelection.LogFilePosition == Beginning
	} else if fromIndex < 0 {
		toIndex += -fromIndex
		fromIndex = 0
		lastPage = logSelection.LogFilePosition == End
	}

	newSelection := Selection{
		ReferencePoint:  *self.createLogLineId(len(self) / 2),
		OffsetFrom:      fromIndex - len(self)/2,
		OffsetTo:        toIndex - len(self)/2,
		LogFilePosition: logSelection.LogFilePosition,
	}
	return self[fromIndex:toIndex], self[fromIndex].Timestamp, self[toIndex-1].Timestamp, newSelection, lastPage
}

func (self LogLines) getLineIndex(logLineId *LogLineId) int {
	if logLineId == nil || logLineId.LogTimestamp == NewestTimestamp || len(self) == 0 || logLineId.LogTimestamp == "" {
		return len(self) - 1
	} else if logLineId.LogTimestamp == OldestTimestamp {
		return 0
	}
	logTimestamp := logLineId.LogTimestamp

	matchingStartedAt := 0
	matchingStartedAt = sort.Search(len(self), func(i int) bool {
		return self[i].Timestamp >= logTimestamp
	})

	linesMatched := 0
	if matchingStartedAt < len(self) && self[matchingStartedAt].Timestamp == logTimestamp {
		for (matchingStartedAt+linesMatched) < len(self) && self[matchingStartedAt+linesMatched].Timestamp == logTimestamp {
			linesMatched += 1
		}
	}

	var offset int
	if logLineId.LineNum < 0 {
		offset = linesMatched + logLineId.LineNum
	} else {
		offset = logLineId.LineNum - 1
	}
	if 0 <= offset && offset < linesMatched {
		return matchingStartedAt + offset
	}
	return LineIndexNotFound
}

func (self LogLines) createLogLineId(lineIndex int) *LogLineId {
	logTimestamp := self[lineIndex].Timestamp
	var step int
	if self[len(self)-1].Timestamp == logTimestamp {
		step = 1
	} else {
		step = -1
	}
	offset := step
	for ; 0 <= lineIndex-offset && lineIndex-offset < len(self); offset += step {
		if self[lineIndex-offset].Timestamp != logTimestamp {
			break
		}
	}
	return &LogLineId{
		LogTimestamp: logTimestamp,
		LineNum:      offset,
	}
}

func ToLogLines(rawLogs string) LogLines {
	logLines := LogLines{}
	for _, line := range strings.Split(rawLogs, "\n") {
		if line != "" {
			startsWithDate := ('0' <= line[0] && line[0] <= '9')
			idx := strings.Index(line, " ")
			if idx > 0 && startsWithDate {
				timestamp := LogTimestamp(line[0:idx])
				content := line[idx+1:]
				logLines = append(logLines, LogLine{Timestamp: timestamp, Content: content})
			} else {
				logLines = append(logLines, LogLine{Timestamp: LogTimestamp("0"), Content: line})
			}
		}
	}
	return logLines
}
