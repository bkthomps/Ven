package search

import "github.com/bkthomps/Ven/buffer"

// Backward Nondeterministic Dawg Matching algorithm

type Pattern struct {
	table  map[rune]uint32
	state  uint32
	length int
}

type Match struct {
	Line        *buffer.Line
	StartOffset int
}

// Compile creates a pattern from a search string
func Compile(search []rune) (pattern *Pattern) {
	pattern = &Pattern{}
	pattern.length = len(search)
	pattern.table = make(map[rune]uint32, 0)
	var x uint32 = 1
	for _, r := range search {
		pattern.table[r] |= x
		x <<= 1
	}
	pattern.state = x - 1
	return pattern
}

// Search starts at the current line, and returns the first match
// that satisfies the pattern, or nil if there are no matches
func (pattern *Pattern) Search(start *buffer.Line) *Match {
	pi := -1
	p := pi + pattern.length
	he := len(start.Data)
	for p < he {
		skip := p
		d := pattern.state
		for d != 0 {
			d &= pattern.table[start.Data[p]]
			p--
			if d == 0 {
				break
			}
			if d&1 != 0 {
				if p != pi {
					skip = p
				} else {
					return &Match{
						Line:        start,
						StartOffset: p + 1,
					}
				}
			}
			d >>= 1
		}
		pi = skip
		p = pi + pattern.length
	}
	return nil
}
