// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vector

// TODO: add tests for NaN and Inf coordinates.

import (
	"image"
	"image/draw"
	"image/png"
	"os"
	"testing"

	"golang.org/x/image/math/f32"
)

// encodePNG is useful for manually debugging the tests.
func encodePNG(dstFilename string, src image.Image) error {
	f, err := os.Create(dstFilename)
	if err != nil {
		return err
	}
	encErr := png.Encode(f, src)
	closeErr := f.Close()
	if encErr != nil {
		return encErr
	}
	return closeErr
}

func TestBasicPath(t *testing.T) {
	mask := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xe3, 0xaa, 0x3e, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfa, 0x5f, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfc, 0x24, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xa1, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfc, 0x14, 0x00, 0x00,
		0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x4a, 0x00, 0x00,
		0x00, 0x00, 0xcc, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x81, 0x00, 0x00,
		0x00, 0x00, 0x66, 0xff, 0xff, 0xff, 0xff, 0xff, 0xef, 0xe4, 0xff, 0xff, 0xff, 0xb6, 0x00, 0x00,
		0x00, 0x00, 0x0c, 0xf2, 0xff, 0xff, 0xfe, 0x9e, 0x15, 0x00, 0x15, 0x96, 0xff, 0xce, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x88, 0xfc, 0xe3, 0x43, 0x00, 0x00, 0x00, 0x00, 0x06, 0xcd, 0xdc, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x10, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x25, 0xde, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x56, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	for _, background := range []uint8{0x00, 0x80} {
		for _, op := range []draw.Op{draw.Over, draw.Src} {
			z := NewRasterizer(16, 16)
			z.MoveTo(f32.Vec2{2, 2})
			z.LineTo(f32.Vec2{8, 2})
			z.QuadTo(f32.Vec2{14, 2}, f32.Vec2{14, 14})
			z.CubeTo(f32.Vec2{8, 2}, f32.Vec2{5, 20}, f32.Vec2{2, 8})
			z.ClosePath()

			dst := image.NewAlpha(z.Bounds())
			for i := range dst.Pix {
				dst.Pix[i] = background
			}
			z.DrawOp = op
			z.Draw(dst, dst.Bounds(), image.Opaque, image.Point{})
			got := dst.Pix

			want := make([]byte, len(mask))
			if op == draw.Over && background == 0x80 {
				for i, ma := range mask {
					want[i] = 0xff - (0xff-ma)/2
				}
			} else {
				copy(want, mask)
			}

			if len(got) != len(want) {
				t.Errorf("background=%#02x, op=%v: len(got)=%d and len(want)=%d differ",
					background, op, len(got), len(want))
				continue
			}
			for i := range got {
				delta := int(got[i]) - int(want[i])
				// The +/- 2 allows different implementations to give different
				// rounding errors.
				if delta < -2 || +2 < delta {
					t.Errorf("background=%#02x, op=%v: i=%d: got %#02x, want %#02x",
						background, op, i, got[i], want[i])
				}
			}
		}
	}
}