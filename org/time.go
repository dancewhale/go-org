package org

import (
	"regexp"
	"time"
)

type TaskTime struct {
	Scheduled string
	Deadline  string
}

var taskTimeRegexp = regexp.MustCompile(`(?i)^\s*((SCHEDULED|DEADLINE):\s*<\d{4}-\d{2}-\d{2}\s+\w+>\s*)+$`)

var scheduledRegexp = regexp.MustCompile(`(?i)^\s*SCHEDULED:\s*<(\d{4}-\d{2}-\d{2}\s+\w+)>\s*`)
var deadlineRegexp = regexp.MustCompile(`(?i)^\s*DEADLINE:\s*<(\d{4}-\d{2}-\d{2}\s+\w+)>\s*`)

func lexTaskTime(line string) (token, bool) {
	if m := taskTimeRegexp.FindStringSubmatch(line); m != nil {
		return token{"taskTime", len(m[1]), m[0], m}, true
	}
	return nilToken, false
}

func (d *Document) parseTaskTime(i int, parentStop stopFn) (int, Node) {
	timeContent := d.tokens[i].content
	i++
	taskTime := TaskTime{}

	if m := scheduledRegexp.FindStringSubmatch(timeContent); m != nil {
		taskTime.Scheduled = m[1]
	}
	if m := deadlineRegexp.FindStringSubmatch(timeContent); m != nil {
		taskTime.Deadline = m[1]
	}
	return 1, taskTime

}

func (n TaskTime) String() string { return String(n) }

func (n TaskTime) GetScheduled() *time.Time {
	if n.Scheduled == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02 Mon", n.Scheduled)
	if err != nil {
		return nil
	} else {
		return &t
	}
}

func (n TaskTime) GetDeadline() *time.Time {
	if n.Deadline == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02 Mon", n.Deadline)
	if err != nil {
		return nil
	} else {
		return &t
	}
}
