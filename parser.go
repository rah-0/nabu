package nabu

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"sort"
	"strings"
	"time"
)

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) FromReader(r io.Reader) *Parser {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		p.lines = append(p.lines, scanner.Text())
	}
	return p
}

func (p *Parser) FromFile(path string) (*Parser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return p.FromReader(f), nil
}

func (p *Parser) FromString(content string) *Parser {
	p.lines = strings.Split(content, "\n")
	return p
}

func (p *Parser) FromLines(lines []string) *Parser {
	p.lines = lines
	return p
}

func (p *Parser) AfterDate(threshold time.Time) *Parser {
	p.afterDate = &threshold
	return p
}

func (p *Parser) Parse() ParsedLogs {
	var parsed ParsedLogs
	traceMap := make(map[string][]Output)
	traceError := make(map[string]string)

	for _, line := range p.lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var entry Output
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		if p.afterDate != nil {
			t, err := time.Parse(TimeLayout, entry.Date)
			if err != nil || !t.After(*p.afterDate) {
				continue
			}
		}
		if entry.UUID != "" {
			if traceError[entry.UUID] == "" && entry.Error != "" {
				traceError[entry.UUID] = entry.Error
			}
			entry.Error = "" // Clear after saving
			traceMap[entry.UUID] = append(traceMap[entry.UUID], entry)
		} else {
			parsed.Entries = append(parsed.Entries, entry)
		}
	}

	for uuid, frames := range traceMap {
		sort.Slice(frames, func(i, j int) bool {
			t1, _ := time.Parse(TimeLayout, frames[i].Date)
			t2, _ := time.Parse(TimeLayout, frames[j].Date)
			return t1.Before(t2)
		})
		parsed.Traces = append(parsed.Traces, ParsedErrorTrace{
			UUID:   uuid,
			Error:  traceError[uuid],
			Frames: frames,
		})
	}
	return parsed
}
