package core

import "math/rand/v2"

var allPieceTypes = [7]PieceType{PieceI, PieceO, PieceT, PieceS, PieceZ, PieceJ, PieceL}

type Bag struct {
	buf  []PieceType
	head int
}

func NewBag() *Bag {
	b := &Bag{}
	b.push()
	b.push()
	return b
}

func (b *Bag) push() {
	set := allPieceTypes
	rand.Shuffle(len(set), func(i, j int) {
		set[i], set[j] = set[j], set[i]
	})
	b.buf = append(b.buf, set[:]...)
}

func (b *Bag) available() int {
	return len(b.buf) - b.head
}

func (b *Bag) Next() PieceType {
	if b.available() < 7 {
		b.push()
	}

	p := b.buf[b.head]
	b.head++

	if b.head >= 7 {
		b.buf = append(b.buf[:0], b.buf[b.head:]...)
		b.head = 0
	}

	return p
}

func (b *Bag) Peek(n int) []PieceType {
	for b.available() < n {
		b.push()
	}
	out := make([]PieceType, n)
	copy(out, b.buf[b.head:b.head+n])
	return out
}
