package tz

import (
	"math/rand"
	"testing"
	"time"

	"github.com/TrueCloudLab/tzhash/gf127"
	"github.com/stretchr/testify/require"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func random() (a *sl2) {
	a = new(sl2)
	a[0][0] = *gf127.Random()
	a[0][1] = *gf127.Random()
	a[1][0] = *gf127.Random()

	// so that result is in SL2
	// d = a^-1*(1+b*c)
	gf127.Mul(&a[0][1], &a[1][0], &a[1][1])
	gf127.Add(&a[1][1], gf127.New(1, 0), &a[1][1])

	t := gf127.New(0, 0)
	gf127.Inv(&a[0][0], t)
	gf127.Mul(t, &a[1][1], &a[1][1])

	return
}

func TestSL2_MarshalBinary(t *testing.T) {
	var (
		a = random()
		b = new(sl2)
	)

	data, err := a.MarshalBinary()
	require.NoError(t, err)

	err = b.UnmarshalBinary(data)
	require.NoError(t, err)

	require.Equal(t, a, b)
}

func TestInv(t *testing.T) {
	var a, b, c *sl2

	c = new(sl2)
	for i := 0; i < 5; i++ {
		a = random()
		b = Inv(a)
		c = c.Mul(a, b)

		require.Equal(t, id, *c)
	}
}
