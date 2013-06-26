// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package adler32 implements the Adler-32 checksum.
//
// It is defined in RFC 1950:
//	Adler-32 is composed of two sums accumulated per byte: s1 is
//	the sum of all bytes, s2 is the sum of all s1 values. Both sums
//	are done modulo 65521. s1 is initialized to 1, s2 to zero.  The
//	Adler-32 checksum is stored as s2*65536 + s1 in most-
//	significant-byte first (network) order.
package rollsum

import "hash"
import "hash/adler32"

const (
	// mod is the largest prime that is less than 65536.
	mod = 65521
	// nmax is the largest n such that
	// 255 * n * (n+1) / 2 + (n+1) * (mod-1) <= 2^32-1.
	// It is mentioned in RFC 1950 (search for "5552").
	nmax = 5552
)

// The size of an Adler-32 checksum in bytes.
const Size = 4

// Rollsum represents the partial evaluation of a checksum.
// The low 16 bits are s1, the high 16 bits are s2.
type Rollsum struct {
	sum uint32
	pos int
	window []byte
	loop bool
}

func (d *Rollsum) Reset() {
	d.sum = 1
	d.pos = 0
	d.loop = false
	for i := 0; i< len(d.window); i++ {
		d.window[i] = 0
	}
}

// New returns a new hash.Hash32 computing the Adler-32 checksum.
func New(n uint32) hash.Hash32 {
	d := new(Rollsum)
	d.window = make([]byte, n)
	d.Reset()
	return d
}

func (d *Rollsum) Size() int { return Size }

func (d *Rollsum) BlockSize() int { return 1 }

func roll(d *Rollsum, p byte) *Rollsum {
	s1, s2 := uint32(d.sum & 0xffff), uint32(d.sum>>16);
	s1 += uint32(p) + uint32((255*len(d.window))/mod + 1)*mod - uint32(d.window[d.pos])
	s2 += s1 - (uint32((len(d.window)) * int(d.window[d.pos])))
	if d.loop {
		s2--;
	}

	s1 %= mod;
	s2 %= mod;

	d.sum = s1 | (s2 << 16);
	d.window[d.pos] = p;
	d.pos = (d.pos + 1) % len(d.window);
	//fmt.Println(d.pos, s1, s2);
	if d.pos == 0 {
		d.loop = true;
	}
	return d;
}

// Add p to the running checksum d.
func update(d *Rollsum, p []byte) *Rollsum {
	for _, x := range p {
		roll(d, x)
	}
	return d
}

func (d *Rollsum) Write(p []byte) (nn int, err error) {
	d = update(d, p)
	return len(p), nil
}

func (d *Rollsum) Sum32() uint32 { return d.sum }

func (d *Rollsum) Sum(in []byte) []byte {
	s := d.Sum32()
	return append(in, byte(s>>24), byte(s>>16), byte(s>>8), byte(s))
}

// Checksum returns the Adler-32 checksum of data.
func Checksum(data []byte) uint32 { return adler32.Checksum(data) }

