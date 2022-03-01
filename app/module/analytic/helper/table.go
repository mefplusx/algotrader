package helper

import (
	"fmt"
	"strings"
)

type Table struct {
	header    []string
	bodyItems [][]string
	colors    []string
}

func (t *Table) SetHeader(items []string) {
	t.header = items
}

func extendBy(target string, bullet string, cnt int) string {
	for i := 0; i < cnt; i++ {
		target += bullet
	}
	return target
}

func (t *Table) Add(args ...interface{}) *Table {
	line := []string{}

	color := ""
	for _, item := range args {
		s := fmt.Sprintf("%v", item)
		if strings.Index(s, "color=") == 0 {
			color = strings.ReplaceAll(s, "color=", "")
		} else {
			line = append(line, s)
		}
	}

	t.bodyItems = append(t.bodyItems, line)
	t.colors = append(t.colors, color)

	return t
}

func (t *Table) prepare(separator string) []string {
	if len(t.bodyItems) == 0 {
		return []string{}
	}
	sizes := make([]int, len(t.bodyItems[0]))

	for i, item := range t.header {
		if len(item) > sizes[i] {
			sizes[i] = len(item)
		}
	}

	for _, scope := range t.bodyItems {
		for i, item := range scope {
			if len(item) > sizes[i] {
				sizes[i] = len(item)
			}
		}
	}

	line := ""
	for i, item := range t.header {
		if line != "" {
			line += separator
		}
		line += extendBy(item, " ", sizes[i]-len(item))
	}
	lines := []string{line}

	for _, scope := range t.bodyItems {
		line = ""
		for i, item := range scope {
			if line != "" {
				line += separator
			}
			line += extendBy(item, " ", sizes[i]-len(item))
		}
		lines = append(lines, line)
	}

	return lines
}

func getCodeColor(colorName string) string {
	switch colorName {
	case "black":
		return "\033[30m"
	case "red":
		return "\033[31m"
	case "green":
		return "\033[32m"
	case "yellow":
		return "\033[33m"
	case "blue":
		return "\033[34m"
	case "magenta":
		return "\033[35m"
	case "cyan":
		return "\033[36m"
	case "white":
		return "\033[37m"
	}

	return ""
}

func (t *Table) Print() {
	const RESET = "\033[0m"
	lines := t.prepare(" ") //t.prepare(" | ")
	for i, line := range lines {
		color := ""
		if i > 0 {
			color = getCodeColor(t.colors[i-1])
		}
		fmt.Println(color + " " + line + " " + RESET)
	}
}
