package search

import (
	"github.com/bkthomps/Ven/buffer"
	"regexp"
)

type Match struct {
	Line        *buffer.Line
	StartOffset int
	Length      int
}

func AllMatches(pattern string, start *buffer.Line) []Match {
	matches := make([]Match, 0)
	re := regexp.MustCompile(pattern)
	for traverse := start; traverse != nil; traverse = traverse.Next {
		indices := re.FindAllStringIndex(string(traverse.Data), -1)
		for _, pair := range indices {
			match := Match{
				Line:        traverse,
				StartOffset: pair[0],
				Length:      pair[1] - pair[0],
			}
			matches = append(matches, match)
		}
	}
	return matches
}
