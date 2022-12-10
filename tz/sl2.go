package tz

import (
	"errors"

	"github.com/TrueCloudLab/tzhash/gf127"
)

type (
	GF127 = gf127.GF127

	sl2 [2][2]GF127
)

var id = sl2{
	{GF127{1, 0}, GF127{0, 0}},
	{GF127{0, 0}, GF127{1, 0}},
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (c *sl2) MarshalBinary() (data []byte, err error) {
	s := c.Bytes()
	return s[:], nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (c *sl2) UnmarshalBinary(data []byte) (err error) {
	if len(data) != 64 {
		return errors.New("data must be 64-bytes long")
	}

	if err = c[0][0].UnmarshalBinary(data[:16]); err != nil {
		return
	}
	if err = c[0][1].UnmarshalBinary(data[16:32]); err != nil {
		return
	}
	if err = c[1][0].UnmarshalBinary(data[32:48]); err != nil {
		return
	}
	if err = c[1][1].UnmarshalBinary(data[48:64]); err != nil {
		return
	}

	return
}

func (c *sl2) mulStrassen(a, b *sl2, x *[8]GF127) *sl2 { //nolint:unused
	// strassen algorithm
	gf127.Add(&a[0][0], &a[1][1], &x[0])
	gf127.Add(&b[0][0], &b[1][1], &x[1])
	gf127.Mul(&x[0], &x[1], &x[0])

	gf127.Add(&a[1][0], &a[1][1], &x[1])
	gf127.Mul(&x[1], &b[0][0], &x[1])

	gf127.Add(&b[0][1], &b[1][1], &x[2])
	gf127.Mul(&x[2], &a[0][0], &x[2])

	gf127.Add(&b[1][0], &b[0][0], &x[3])
	gf127.Mul(&x[3], &a[1][1], &x[3])

	gf127.Add(&a[0][0], &a[0][1], &x[4])
	gf127.Mul(&x[4], &b[1][1], &x[4])

	gf127.Add(&a[1][0], &a[0][0], &x[5])
	gf127.Add(&b[0][0], &b[0][1], &x[6])
	gf127.Mul(&x[5], &x[6], &x[5])

	gf127.Add(&a[0][1], &a[1][1], &x[6])
	gf127.Add(&b[1][0], &b[1][1], &x[7])
	gf127.Mul(&x[6], &x[7], &x[6])

	gf127.Add(&x[2], &x[4], &c[0][1])
	gf127.Add(&x[1], &x[3], &c[1][0])

	gf127.Add(&x[4], &x[6], &x[4])
	gf127.Add(&x[0], &x[3], &c[0][0])
	gf127.Add(&c[0][0], &x[4], &c[0][0])

	gf127.Add(&x[0], &x[1], &x[0])
	gf127.Add(&x[2], &x[5], &c[1][1])
	gf127.Add(&c[1][1], &x[0], &c[1][1])

	return c
}

func (c *sl2) MulA() *sl2 {
	var a GF127

	gf127.Mul10(&c[0][0], &a)
	gf127.Mul1(&c[0][0], &c[0][1])
	gf127.Add(&a, &c[0][1], &c[0][0])

	gf127.Mul10(&c[1][0], &a)
	gf127.Mul1(&c[1][0], &c[1][1])
	gf127.Add(&a, &c[1][1], &c[1][0])

	return c
}

func (c *sl2) MulB() *sl2 {
	var a GF127

	gf127.Mul1(&c[0][0], &a)
	gf127.Mul10(&c[0][0], &c[0][0])
	gf127.Add(&c[0][1], &c[0][0], &c[0][0])
	gf127.Add(&c[0][0], &a, &c[0][1])

	gf127.Mul1(&c[1][0], &a)
	gf127.Mul10(&c[1][0], &c[1][0])
	gf127.Add(&c[1][1], &c[1][0], &c[1][0])
	gf127.Add(&c[1][0], &a, &c[1][1])

	return c
}

// Mul returns a * b in GL_2(GF(2^127))
func (c *sl2) Mul(a, b *sl2) *sl2 {
	var x [4]GF127

	gf127.Mul(&a[0][0], &b[0][0], &x[0])
	gf127.Mul(&a[0][0], &b[0][1], &x[1])
	gf127.Mul(&a[1][0], &b[0][0], &x[2])
	gf127.Mul(&a[1][0], &b[0][1], &x[3])

	gf127.Mul(&a[0][1], &b[1][0], &c[0][0])
	gf127.Add(&c[0][0], &x[0], &c[0][0])
	gf127.Mul(&a[0][1], &b[1][1], &c[0][1])
	gf127.Add(&c[0][1], &x[1], &c[0][1])
	gf127.Mul(&a[1][1], &b[1][0], &c[1][0])
	gf127.Add(&c[1][0], &x[2], &c[1][0])
	gf127.Mul(&a[1][1], &b[1][1], &c[1][1])
	gf127.Add(&c[1][1], &x[3], &c[1][1])
	return c
}

// Inv returns inverse of a in GL_2(GF(2^127))
func Inv(a *sl2) (b *sl2) {
	b = new(sl2)
	inv(a, b, new([2]GF127))
	return
}

func inv(a, b *sl2, t *[2]GF127) {
	gf127.Mul(&a[0][0], &a[1][1], &t[0])
	gf127.Mul(&a[0][1], &a[1][0], &t[1])
	gf127.Add(&t[0], &t[1], &t[0])
	gf127.Inv(&t[0], &t[1])

	gf127.Mul(&t[1], &a[0][0], &b[1][1])
	gf127.Mul(&t[1], &a[0][1], &b[0][1])
	gf127.Mul(&t[1], &a[1][0], &b[1][0])
	gf127.Mul(&t[1], &a[1][1], &b[0][0])
}

func (c *sl2) String() string {
	return c[0][0].String() + c[0][1].String() +
		c[1][0].String() + c[1][1].String()
}

func (c *sl2) Bytes() (b [Size]byte) {
	t := c[0][0].Bytes()
	copy(b[:], t[:])

	t = c[0][1].Bytes()
	copy(b[16:], t[:])

	t = c[1][0].Bytes()
	copy(b[32:], t[:])

	t = c[1][1].Bytes()
	copy(b[48:], t[:])

	return
}
