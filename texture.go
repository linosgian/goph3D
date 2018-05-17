package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/disintegration/imaging"
	"github.com/go-gl/gl/v4.3-core/gl"
)

type Texture struct {
	rendererID  uint32
	filepath    string
	localBuffer []uint8
	Width       int32
	Height      int32
	bpp         int
}

func NewTexture(filepath string) (*Texture, error) {
	im, err := ReadImageFile(filepath)
	if err != nil {
		return nil, err
	}

	t := Texture{
		localBuffer: im.Pix,
		Width:       int32(im.Rect.Size().X),
		Height:      int32(im.Rect.Size().Y),
		filepath:    filepath,
	}

	gl.GenTextures(1, &t.rendererID)
	gl.BindTexture(gl.TEXTURE_2D, t.rendererID)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, t.Width, t.Height, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(t.localBuffer))
	gl.BindTexture(gl.TEXTURE_2D, 0)
	return &t, nil
}

func (t *Texture) Delete() {
	gl.DeleteTextures(1, &t.rendererID)
}

func (t *Texture) Bind(slot uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + slot)
	gl.BindTexture(gl.TEXTURE_2D, t.rendererID)
}

func (t *Texture) Unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func ReadImageFile(filepath string) (*image.NRGBA, error) {
	r, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not open image file: %v\n", err)
	}
	defer r.Close()

	im, _, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("could not decode image file: %v\n", err)
	}
	rgba := image.NewRGBA(im.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), im, image.Point{0, 0}, draw.Src)
	// NOTE: Instead of rotating, the image itself should be rotated
	// already when on disk.
	return imaging.Rotate180(rgba), nil
}
