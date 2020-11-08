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
		strData := string(traverse.Data)
		indices := re.FindAllStringIndex(strData, -1)
		if len(indices) == 0 {
			continue
		}
		if count == 0 {
			count++
		}
		byteToRuneIndex := make(map[int]int, 0)
		runeIndex := 0
		for byteIndex := range strData {
			byteToRuneIndex[byteIndex] = runeIndex
			runeIndex++
		}
		matchInstances := make([]MatchInstance, 0)
		for _, pair := range indices {
			instance := MatchInstance{
				StartOffset: byteToRuneIndex[pair[0]],
				Length:      byteToRuneIndex[pair[1]] - byteToRuneIndex[pair[0]],
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
