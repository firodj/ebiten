// Copyright 2020 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ebiten_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestShaderFill(t *testing.T) {
	const w, h = 16, 16

	dst := ebiten.NewImage(w, h)
	s, err := ebiten.NewShader([]byte(`package main

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	return vec4(1, 0, 0, 1)
}
`))
	if err != nil {
		t.Fatal(err)
	}

	dst.DrawRectShader(w/2, h/2, s, nil)

	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			got := dst.At(i, j).(color.RGBA)
			var want color.RGBA
			if i < w/2 && j < h/2 {
				want = color.RGBA{0xff, 0, 0, 0xff}
			}
			if got != want {
				t.Errorf("dst.At(%d, %d): got: %v, want: %v", i, j, got, want)
			}
		}
	}
}

func TestShaderFillWithDrawImage(t *testing.T) {
	const w, h = 16, 16

	dst := ebiten.NewImage(w, h)
	s, err := ebiten.NewShader([]byte(`package main

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	return vec4(1, 0, 0, 1)
}
`))
	if err != nil {
		t.Fatal(err)
	}

	src := ebiten.NewImage(w/2, h/2)
	op := &ebiten.DrawRectShaderOptions{}
	op.Images[0] = src
	dst.DrawRectShader(w/2, h/2, s, op)

	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			got := dst.At(i, j).(color.RGBA)
			var want color.RGBA
			if i < w/2 && j < h/2 {
				want = color.RGBA{0xff, 0, 0, 0xff}
			}
			if got != want {
				t.Errorf("dst.At(%d, %d): got: %v, want: %v", i, j, got, want)
			}
		}
	}
}

func TestShaderFillWithDrawTriangles(t *testing.T) {
	const w, h = 16, 16

	dst := ebiten.NewImage(w, h)
	s, err := ebiten.NewShader([]byte(`package main

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	return vec4(1, 0, 0, 1)
}
`))
	if err != nil {
		t.Fatal(err)
	}

	src := ebiten.NewImage(w/2, h/2)
	op := &ebiten.DrawTrianglesShaderOptions{}
	op.Images[0] = src

	vs := []ebiten.Vertex{
		{
			DstX:   0,
			DstY:   0,
			SrcX:   0,
			SrcY:   0,
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		},
		{
			DstX:   w,
			DstY:   0,
			SrcX:   w / 2,
			SrcY:   0,
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		},
		{
			DstX:   0,
			DstY:   h,
			SrcX:   0,
			SrcY:   h / 2,
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		},
		{
			DstX:   w,
			DstY:   h,
			SrcX:   w / 2,
			SrcY:   h / 2,
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		},
	}
	is := []uint16{0, 1, 2, 1, 2, 3}

	dst.DrawTrianglesShader(vs, is, s, op)

	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			got := dst.At(i, j).(color.RGBA)
			want := color.RGBA{0xff, 0, 0, 0xff}
			if got != want {
				t.Errorf("dst.At(%d, %d): got: %v, want: %v", i, j, got, want)
			}
		}
	}
}

func TestShaderFunction(t *testing.T) {
	const w, h = 16, 16

	dst := ebiten.NewImage(w, h)
	s, err := ebiten.NewShader([]byte(`package main

func clr(red float) (float, float, float, float) {
	return red, 0, 0, 1
}

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	return vec4(clr(1))
}
`))
	if err != nil {
		t.Fatal(err)
	}

	dst.DrawRectShader(w, h, s, nil)

	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			got := dst.At(i, j).(color.RGBA)
			want := color.RGBA{0xff, 0, 0, 0xff}
			if got != want {
				t.Errorf("dst.At(%d, %d): got: %v, want: %v", i, j, got, want)
			}
		}
	}
}

func TestShaderUninitializedUniformVariables(t *testing.T) {
	const w, h = 16, 16

	dst := ebiten.NewImage(w, h)
	s, err := ebiten.NewShader([]byte(`package main

var U vec4

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	return U
}
`))
	if err != nil {
		t.Fatal(err)
	}

	dst.DrawRectShader(w, h, s, nil)

	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			got := dst.At(i, j).(color.RGBA)
			var want color.RGBA
			if got != want {
				t.Errorf("dst.At(%d, %d): got: %v, want: %v", i, j, got, want)
			}
		}
	}
}

func TestShaderMatrix(t *testing.T) {
	const w, h = 16, 16

	dst := ebiten.NewImage(w, h)
	s, err := ebiten.NewShader([]byte(`package main

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	var a, b mat4
	a[0] = vec4(0.125, 0.0625, 0.0625, 0.0625)
	a[1] = vec4(0.25, 0.25, 0.0625, 0.1875)
	a[2] = vec4(0.1875, 0.125, 0.25, 0.25)
	a[3] = vec4(0.0625, 0.1875, 0.125, 0.25)
	b[0] = vec4(0.0625, 0.125, 0.0625, 0.125)
	b[1] = vec4(0.125, 0.1875, 0.25, 0.0625)
	b[2] = vec4(0.125, 0.125, 0.1875, 0.1875)
	b[3] = vec4(0.25, 0.0625, 0.125, 0.0625)
	return vec4((a * b * vec4(1, 1, 1, 1)).xyz, 1)
}
`))
	if err != nil {
		t.Fatal(err)
	}

	src := ebiten.NewImage(w, h)
	op := &ebiten.DrawRectShaderOptions{}
	op.Images[0] = src
	dst.DrawRectShader(w, h, s, op)

	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			got := dst.At(i, j).(color.RGBA)
			want := color.RGBA{87, 82, 71, 255}
			if got != want {
				t.Errorf("dst.At(%d, %d): got: %v, want: %v", i, j, got, want)
			}
		}
	}
}

func TestShaderSubImage(t *testing.T) {
	const w, h = 16, 16

	s, err := ebiten.NewShader([]byte(`package main

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	r := imageSrc0At(texCoord).r
	g := imageSrc1At(texCoord).g
	return vec4(r, g, 0, 1)
}
`))
	if err != nil {
		t.Fatal(err)
	}

	src0 := ebiten.NewImage(w, h)
	pix0 := make([]byte, 4*w*h)
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			if 2 <= i && i < 10 && 3 <= j && j < 11 {
				pix0[4*(j*w+i)] = 0xff
				pix0[4*(j*w+i)+1] = 0
				pix0[4*(j*w+i)+2] = 0
				pix0[4*(j*w+i)+3] = 0xff
			}
		}
	}
	src0.ReplacePixels(pix0)
	src0 = src0.SubImage(image.Rect(2, 3, 10, 11)).(*ebiten.Image)

	src1 := ebiten.NewImage(w, h)
	pix1 := make([]byte, 4*w*h)
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			if 6 <= i && i < 14 && 8 <= j && j < 16 {
				pix1[4*(j*w+i)] = 0
				pix1[4*(j*w+i)+1] = 0xff
				pix1[4*(j*w+i)+2] = 0
				pix1[4*(j*w+i)+3] = 0xff
			}
		}
	}
	src1.ReplacePixels(pix1)
	src1 = src1.SubImage(image.Rect(6, 8, 14, 16)).(*ebiten.Image)

	testPixels := func(testname string, dst *ebiten.Image) {
		for j := 0; j < h; j++ {
			for i := 0; i < w; i++ {
				got := dst.At(i, j).(color.RGBA)
				var want color.RGBA
				if i < w/2 && j < h/2 {
					want = color.RGBA{0xff, 0xff, 0, 0xff}
				}
				if got != want {
					t.Errorf("%s dst.At(%d, %d): got: %v, want: %v", testname, i, j, got, want)
				}
			}
		}
	}

	t.Run("DrawRectShader", func(t *testing.T) {
		dst := ebiten.NewImage(w, h)
		op := &ebiten.DrawRectShaderOptions{}
		op.Images[0] = src0
		op.Images[1] = src1
		dst.DrawRectShader(w/2, h/2, s, op)
		testPixels("DrawRectShader", dst)
	})

	t.Run("DrawTrianglesShader", func(t *testing.T) {
		dst := ebiten.NewImage(w, h)
		vs := []ebiten.Vertex{
			{
				DstX:   0,
				DstY:   0,
				SrcX:   2,
				SrcY:   3,
				ColorR: 1,
				ColorG: 1,
				ColorB: 1,
				ColorA: 1,
			},
			{
				DstX:   w / 2,
				DstY:   0,
				SrcX:   10,
				SrcY:   3,
				ColorR: 1,
				ColorG: 1,
				ColorB: 1,
				ColorA: 1,
			},
			{
				DstX:   0,
				DstY:   h / 2,
				SrcX:   2,
				SrcY:   11,
				ColorR: 1,
				ColorG: 1,
				ColorB: 1,
				ColorA: 1,
			},
			{
				DstX:   w / 2,
				DstY:   h / 2,
				SrcX:   10,
				SrcY:   11,
				ColorR: 1,
				ColorG: 1,
				ColorB: 1,
				ColorA: 1,
			},
		}
		is := []uint16{0, 1, 2, 1, 2, 3}

		op := &ebiten.DrawTrianglesShaderOptions{}
		op.Images[0] = src0
		op.Images[1] = src1
		dst.DrawTrianglesShader(vs, is, s, op)
		testPixels("DrawTrianglesShader", dst)
	})
}

// Issue #1404
func TestShaderDerivatives(t *testing.T) {
	const w, h = 16, 16

	s, err := ebiten.NewShader([]byte(`package main

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	p := imageSrc0At(texCoord)
	return vec4(abs(dfdx(p.r)), abs(dfdy(p.g)), 0, 1)
}
`))
	if err != nil {
		t.Fatal(err)
	}

	dst := ebiten.NewImage(w, h)
	src := ebiten.NewImage(w, h)
	pix := make([]byte, 4*w*h)
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			if i < w/2 {
				pix[4*(j*w+i)] = 0xff
			}
			if j < h/2 {
				pix[4*(j*w+i)+1] = 0xff
			}
			pix[4*(j*w+i)+3] = 0xff
		}
	}
	src.ReplacePixels(pix)

	op := &ebiten.DrawRectShaderOptions{}
	op.Images[0] = src
	dst.DrawRectShader(w, h, s, op)

	// The results of the edges might be unreliable. Skip the edges.
	for j := 1; j < h-1; j++ {
		for i := 1; i < w-1; i++ {
			got := dst.At(i, j).(color.RGBA)
			want := color.RGBA{0, 0, 0, 0xff}
			if i == w/2-1 || i == w/2 {
				want.R = 0xff
			}
			if j == h/2-1 || j == h/2 {
				want.G = 0xff
			}
			if got != want {
				t.Errorf("dst.At(%d, %d): got: %v, want: %v", i, j, got, want)
			}
		}
	}
}

// Issue #1701
func TestShaderDerivatives2(t *testing.T) {
	const w, h = 16, 16

	s, err := ebiten.NewShader([]byte(`package main

// This function uses dfdx and then should not be in GLSL's vertex shader (#1701).
func Foo(p vec4) vec4 {
	return vec4(abs(dfdx(p.r)), abs(dfdy(p.g)), 0, 1)
}

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	p := imageSrc0At(texCoord)
	return Foo(p)
}
`))
	if err != nil {
		t.Fatal(err)
	}

	dst := ebiten.NewImage(w, h)
	src := ebiten.NewImage(w, h)
	pix := make([]byte, 4*w*h)
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			if i < w/2 {
				pix[4*(j*w+i)] = 0xff
			}
			if j < h/2 {
				pix[4*(j*w+i)+1] = 0xff
			}
			pix[4*(j*w+i)+3] = 0xff
		}
	}
	src.ReplacePixels(pix)

	op := &ebiten.DrawRectShaderOptions{}
	op.Images[0] = src
	dst.DrawRectShader(w, h, s, op)

	// The results of the edges might be unreliable. Skip the edges.
	for j := 1; j < h-1; j++ {
		for i := 1; i < w-1; i++ {
			got := dst.At(i, j).(color.RGBA)
			want := color.RGBA{0, 0, 0, 0xff}
			if i == w/2-1 || i == w/2 {
				want.R = 0xff
			}
			if j == h/2-1 || j == h/2 {
				want.G = 0xff
			}
			if got != want {
				t.Errorf("dst.At(%d, %d): got: %v, want: %v", i, j, got, want)
			}
		}
	}
}

// Issue #1754
func TestShaderUniformFirstElement(t *testing.T) {
	shaders := []struct {
		Name     string
		Shader   string
		Uniforms map[string]interface{}
	}{
		{
			Name: "float array",
			Shader: `package main

var C [2]float

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	return vec4(C[0], 1, 1, 1)
}`,
			Uniforms: map[string]interface{}{
				"C": []float32{1, 1},
			},
		},
		{
			Name: "float one-element array",
			Shader: `package main

var C [1]float

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	return vec4(C[0], 1, 1, 1)
}`,
			Uniforms: map[string]interface{}{
				"C": []float32{1},
			},
		},
		{
			Name: "matrix array",
			Shader: `package main

var C [2]mat2

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	return vec4(C[0][0][0], 1, 1, 1)
}`,
			Uniforms: map[string]interface{}{
				"C": []float32{1, 0, 0, 0, 0, 0, 0, 0},
			},
		},
	}

	for _, shader := range shaders {
		shader := shader
		t.Run(shader.Name, func(t *testing.T) {
			const w, h = 1, 1

			dst := ebiten.NewImage(w, h)
			s, err := ebiten.NewShader([]byte(shader.Shader))
			if err != nil {
				t.Fatal(err)
			}

			op := &ebiten.DrawRectShaderOptions{}
			op.Uniforms = shader.Uniforms
			dst.DrawRectShader(w, h, s, op)
			if got, want := dst.At(0, 0), (color.RGBA{0xff, 0xff, 0xff, 0xff}); got != want {
				t.Errorf("got: %v, want: %v", got, want)
			}
		})
	}
}

// Issue #2006
func TestShaderFuncMod(t *testing.T) {
	const w, h = 16, 16

	dst := ebiten.NewImage(w, h)
	s, err := ebiten.NewShader([]byte(`package main

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	r := mod(-0.25, 1.0)
	return vec4(r, 0, 0, 1)
}
`))
	if err != nil {
		t.Fatal(err)
	}

	dst.DrawRectShader(w/2, h/2, s, nil)

	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			got := dst.At(i, j).(color.RGBA)
			var want color.RGBA
			if i < w/2 && j < h/2 {
				want = color.RGBA{0xc0, 0, 0, 0xff}
			}
			if !sameColors(got, want, 2) {
				t.Errorf("dst.At(%d, %d): got: %v, want: %v", i, j, got, want)
			}
		}
	}
}

func TestShaderMatrixInitialize(t *testing.T) {
	const w, h = 16, 16

	src := ebiten.NewImage(w, h)
	src.Fill(color.RGBA{0x10, 0x20, 0x30, 0xff})

	dst := ebiten.NewImage(w, h)
	s, err := ebiten.NewShader([]byte(`package main

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	return mat4(2) * imageSrc0At(texCoord);
}
`))
	if err != nil {
		t.Fatal(err)
	}

	op := &ebiten.DrawRectShaderOptions{}
	op.Images[0] = src
	dst.DrawRectShader(w, h, s, op)

	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			got := dst.At(i, j).(color.RGBA)
			want := color.RGBA{0x20, 0x40, 0x60, 0xff}
			if !sameColors(got, want, 2) {
				t.Errorf("dst.At(%d, %d): got: %v, want: %v", i, j, got, want)
			}
		}
	}
}
