package search

import (
	"regexp"

	"github.com/bkthomps/Ven/buffer"
)

type MatchLine struct {
	Line      *buffer.Line
	Instances []MatchInstance
}

type MatchInstance struct {
	StartOffset int
	Length      int
}

func AllMatches(pattern string, start *buffer.Line, maxLineCount int) (matches []MatchLine, firstLineIndex int) {
	count := 0
	firstLineIndex = 0
	matches = make([]MatchLine, 0)
	re, err := regexp.Compile(pattern)
	if err != nil {
		return matches, 0
	}
	for traverse := start; traverse != nil; traverse = traverse.Next {
		if count > 0 {
			count++
			if count > maxLineCount {
				return matches, firstLineIndex
			}
		}
		if count == 0 {
			firstLineIndex++
		}
		indices := re.FindAllStringIndex(string(traverse.Data), -1)
		if len(indices) == 0 {
			continue
		}
		if count == 0 {
			count++
		}
		matchInstances := make([]MatchInstance, 0)
		for _, pair := range indices {
			instance := MatchInstance{
				StartOffset: pair[0],
				Length:      pair[1] - pair[0],
			}
			matchInstances = append(matchInstances, instance)
		}
		match := MatchLine{
			Line:      traverse,
			Instances: matchInstances,
		}
		matches = append(matches, match)
	}
	return matches, firstLineIndex
}
