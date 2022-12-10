//go:build amd64 && !generic
// +build amd64,!generic

package tz

import (
	"github.com/TrueCloudLab/tzhash/gf127"
	"golang.org/x/sys/cpu"
)

func write(d *digest, data []byte) (n int, err error) {
	switch {
	case cpu.X86.HasAVX && cpu.X86.HasAVX2:
		return writeAVX2(d, data)
	case cpu.X86.HasAVX:
		return writeAVX(d, data)
	default:
		return writeGeneric(d, data)
	}
}

func writeAVX2(d *digest, data []byte) (n int, err error) {
	n = len(data)
	if len(data) != 0 {
		mulByteSliceRightx2(&d.x[0], &d.x[2], n, &data[0])
	}
	return
}

func writeAVX(d *digest, data []byte) (n int, err error) {
	n = len(data)
	for _, b := range data {
		mulByteRight(&d.x[0], &d.x[1], &d.x[2], &d.x[3], b)
	}
	return
}

func mulByteRight(c00, c01, c10, c11 *GF127, b byte)
func mulByteSliceRightx2(c00c10 *gf127.GF127, c01c11 *gf127.GF127, n int, data *byte)
