package core

import "testing"

func TestBag_EvenDistribution(t *testing.T) {
	b := NewBag()
	counts := make(map[PieceType]int)

	const total = 7000
	for range total {
		counts[b.Next()]++
	}

	expected := total / 7
	for _, pt := range allPieceTypes {
		if counts[pt] != expected {
			t.Errorf("piece %v: got %d, want %d", pt, counts[pt], expected)
		}
	}
}

func TestBag_PeekDoesNotConsume(t *testing.T) {
	b := NewBag()

	peeked := b.Peek(3)
	for i, want := range peeked {
		got := b.Next()
		if got != want {
			t.Errorf("step %d: Peek=%v, Next=%v", i, want, got)
		}
	}
}

func TestBag_PeekIsStable(t *testing.T) {
	b := NewBag()

	first := b.Peek(5)
	second := b.Peek(5)

	for i := range first {
		if first[i] != second[i] {
			t.Errorf("peek[%d] non stable: %v: %v", i, first[i], second[i])
		}
	}
}

func TestBag_NoDuplicatesInBag(t *testing.T) {
	b := NewBag()

	for bag := range 200 {
		seen := make(map[PieceType]bool, 7)
		for range 7 {
			p := b.Next()
			if seen[p] {
				t.Errorf("diplicte %v in set %d", p, bag)
			}
			seen[p] = true
		}
	}
}

func TestBag_PeekBeyondCurrentSet(t *testing.T) {
	b := NewBag()

	for range 6 {
		b.Next()
	}

	peeked := b.Peek(5)
	if len(peeked) != 5 {
		t.Fatalf("Peek(5) retunr %d elements", len(peeked))
	}

	for i, want := range peeked {
		got := b.Next()
		if got != want {
			t.Errorf("step %d: Peek=%v, Next=%v", i, want, got)
		}
	}
}
