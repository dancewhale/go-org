package org

import (
	"regexp"
	"strings"
)

type Drawer struct {
	Name     string
	Children []Node
}

type PropertyDrawer struct {
	Properties [][]string
}

type LogBookDrawer struct {
	Clocks []Clock
}

type Clock struct {
	Start string
	End   string
}

var beginDrawerRegexp = regexp.MustCompile(`^(\s*):(\S+):\s*$`)
var endDrawerRegexp = regexp.MustCompile(`(?i)^(\s*):END:\s*$`)
var propertyRegexp = regexp.MustCompile(`^(\s*):(\S+):(\s+(.*)$|$)`)

var logBookClockRegexp = regexp.MustCompile(`^\s*CLOCK:\s+\[(\d{4}-\d{2}-\d{2}\s+\w{3}\s+\d{2}:\d{2})\]-+\[(\d{4}-\d{2}-\d{2}\s+\w{3}\s+\d{2}:\d{2})\]`)

func lexDrawer(line string) (token, bool) {
	if m := endDrawerRegexp.FindStringSubmatch(line); m != nil {
		return token{"endDrawer", len(m[1]), "", m}, true
	} else if m := beginDrawerRegexp.FindStringSubmatch(line); m != nil {
		return token{"beginDrawer", len(m[1]), strings.ToUpper(m[2]), m}, true
	}
	return nilToken, false
}

func (d *Document) parseDrawer(i int, parentStop stopFn) (int, Node) {
	name := strings.ToUpper(d.tokens[i].content)
	if name == "PROPERTIES" {
		return d.parsePropertyDrawer(i, parentStop)
	} else if name == "LOGBOOK" {
		return d.parseLogBookClock(i, parentStop)
	}
	drawer, start := Drawer{Name: name}, i
	i++
	stop := func(d *Document, i int) bool {
		if parentStop(d, i) {
			return true
		}
		kind := d.tokens[i].kind
		return kind == "beginDrawer" || kind == "endDrawer" || kind == "headline"
	}
	for {
		consumed, nodes := d.parseMany(i, stop)
		i += consumed
		drawer.Children = append(drawer.Children, nodes...)
		if i < len(d.tokens) && d.tokens[i].kind == "beginDrawer" {
			p := Paragraph{[]Node{Text{":" + d.tokens[i].content + ":", false}}}
			drawer.Children = append(drawer.Children, p)
			i++
		} else {
			break
		}
	}
	if i < len(d.tokens) && d.tokens[i].kind == "endDrawer" {
		i++
	}
	return i - start, drawer
}

func (d *Document) parsePropertyDrawer(i int, parentStop stopFn) (int, Node) {
	drawer, start := PropertyDrawer{}, i
	i++
	stop := func(d *Document, i int) bool {
		return parentStop(d, i) || (d.tokens[i].kind != "text" && d.tokens[i].kind != "beginDrawer")
	}
	for ; !stop(d, i); i++ {
		m := propertyRegexp.FindStringSubmatch(d.tokens[i].matches[0])
		if m == nil {
			return 0, nil
		}
		k, v := strings.ToUpper(m[2]), strings.TrimSpace(m[4])
		drawer.Properties = append(drawer.Properties, []string{k, v})
	}
	if i < len(d.tokens) && d.tokens[i].kind == "endDrawer" {
		i++
	} else {
		return 0, nil
	}
	return i - start, drawer
}

func (d *Document) parseLogBookClock(i int, parentStop stopFn) (int, Node) {
	drawer, start := LogBookDrawer{}, i
	i++
	stop := func(d *Document, i int) bool {
		return parentStop(d, i) || (d.tokens[i].kind != "text" && d.tokens[i].kind != "beginDrawer")
	}
	for ; !stop(d, i); i++ {
		m := logBookClockRegexp.FindStringSubmatch(d.tokens[i].matches[0])
		if m == nil {
			return 0, nil
		}
		clock := Clock{Start: m[1], End: m[2]}
		drawer.Clocks = append(drawer.Clocks, clock)
	}
	if i < len(d.tokens) && d.tokens[i].kind == "endDrawer" {
		i++
	} else {
		return 0, nil
	}
	return i - start, drawer
}

func (d *PropertyDrawer) Get(key string) (string, bool) {
	if d == nil {
		return "", false
	}
	for _, kvPair := range d.Properties {
		if kvPair[0] == key {
			return kvPair[1], true
		}
	}
	return "", false
}

func (n Drawer) String() string         { return String(n) }
func (n PropertyDrawer) String() string { return String(n) }
func (n Clock) String() string          { return String(n) }
func (n LogBookDrawer) String() string  { return String(n) }
