package buffer

import "testing"

func TestSpacingOffset(t *testing.T) {
	file := File{}
	file.Init("")
	runes := []rune{'a', '\t', 'b', 'c', ' ', ' ', '\t', '汉', '字', '🦀', '😂'}
	spacing := []int{1, 8, 9, 10, 11, 12, 16, 18, 20, 22, 24}
	for i, r := range runes {
		file.Current.AddAt(i, r)
	}
	file.spacingOffset = 0
	for i, r := range file.Current.Data {
		if r != runes[i] {
			t.Error("incorrect rune data")
		}
		file.spacingOffset = file.runeWidthIncrease(r)
		file.runeOffset++
		if file.spacingOffset != spacing[i] {
			t.Errorf("bad spacing offset: expected %d, received %d", spacing[i], file.spacingOffset)
		}
	}
	for i := len(file.Current.Data) - 1; i >= 0; i-- {
		if file.spacingOffset != spacing[i] {
			t.Errorf("bad spacing offset: expected %d, received %d", spacing[i], file.spacingOffset)
		}
		file.runeOffset--
		file.spacingOffset = file.runeWidthDecrease(file.Current.Data[i])
	}
}
