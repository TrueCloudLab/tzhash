package tz

import (
	"github.com/TrueCloudLab/tzhash/gf127"
)

const (
	// Size is the size of a Tillich-Zémor hash sum in bytes.
	Size          = 64
	hashBlockSize = 128
)

type digest struct {
	// Stores matrix cells in the following order:
	// [ 0 2 ]
	// [ 1 3 ]
	// This is done to reuse the same digest between generic
	// and AVX2 implementation.
	x [4]GF127
}

// New returns a new hash.Hash computing the Tillich-Zémor checksum.
func New() *digest {
	d := new(digest)
	d.Reset()
	return d
}

// Sum returns Tillich-Zémor checksum of data.
func Sum(data []byte) [Size]byte {
	d := new(digest)
	d.Reset()
	_, _ = d.Write(data) // no errors
	return d.checkSum()
}

// Sum implements hash.Hash.
func (d *digest) Sum(in []byte) []byte {
	// Make a copy of d so that caller can keep writing and summing.
	d0 := *d
	h := d0.checkSum()
	return append(in, h[:]...)
}

func (d *digest) checkSum() (b [Size]byte) {
	t := d.x[0].Bytes()
	copy(b[:], t[:])

	t = d.x[2].Bytes()
	copy(b[16:], t[:])

	t = d.x[1].Bytes()
	copy(b[32:], t[:])

	t = d.x[3].Bytes()
	copy(b[48:], t[:])

	return
}

// Reset implements hash.Hash.
func (d *digest) Reset() {
	d.x[0] = GF127{1, 0}
	d.x[1] = GF127{0, 0}
	d.x[2] = GF127{0, 0}
	d.x[3] = GF127{1, 0}
}

// Write implements hash.Hash.
func (d *digest) Write(data []byte) (n int, err error) {
	return write(d, data)
}

func writeGeneric(d *digest, data []byte) (n int, err error) {
	n = len(data)
	tmp := new(GF127)
	for _, b := range data {
		mulBitRightGeneric(&d.x[0], &d.x[1], &d.x[2], &d.x[3], b&0x80 != 0, tmp)
		mulBitRightGeneric(&d.x[0], &d.x[1], &d.x[2], &d.x[3], b&0x40 != 0, tmp)
		mulBitRightGeneric(&d.x[0], &d.x[1], &d.x[2], &d.x[3], b&0x20 != 0, tmp)
		mulBitRightGeneric(&d.x[0], &d.x[1], &d.x[2], &d.x[3], b&0x10 != 0, tmp)
		mulBitRightGeneric(&d.x[0], &d.x[1], &d.x[2], &d.x[3], b&0x08 != 0, tmp)
		mulBitRightGeneric(&d.x[0], &d.x[1], &d.x[2], &d.x[3], b&0x04 != 0, tmp)
		mulBitRightGeneric(&d.x[0], &d.x[1], &d.x[2], &d.x[3], b&0x02 != 0, tmp)
		mulBitRightGeneric(&d.x[0], &d.x[1], &d.x[2], &d.x[3], b&0x01 != 0, tmp)
	}
	return
}

// Size implements hash.Hash.
func (d *digest) Size() int {
	return Size
}

// BlockSize implements hash.Hash.
func (d *digest) BlockSize() int {
	return hashBlockSize
}

func mulBitRightGeneric(c00, c10, c01, c11 *GF127, bit bool, tmp *GF127) {
	if bit {
		*tmp = *c00
		gf127.Mul10(c00, c00)
		gf127.Add(c00, c01, c00)
		gf127.Mul11(tmp, tmp)
		gf127.Add(c01, tmp, c01)

		*tmp = *c10
		gf127.Mul10(c10, c10)
		gf127.Add(c10, c11, c10)
		gf127.Mul11(tmp, tmp)
		gf127.Add(c11, tmp, c11)
	} else {
		*tmp = *c00
		gf127.Mul10(c00, c00)
		gf127.Add(c00, c01, c00)
		*c01 = *tmp

		*tmp = *c10
		gf127.Mul10(c10, c10)
		gf127.Add(c10, c11, c10)
		*c11 = *tmp
	}
}
