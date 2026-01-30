package core

import (
	"bytes"
	_ "embed"
	"image"
	"image/png"
)

//go:embed assets/bg_48.png
var bg48Data []byte

//go:embed assets/bg_96.png
var bg96Data []byte

var (
	Bg48 image.Image
	Bg96 image.Image
)

func init() {
	var err error
	Bg48, err = png.Decode(bytes.NewReader(bg48Data))
	if err != nil {
		panic(err)
	}
	Bg96, err = png.Decode(bytes.NewReader(bg96Data))
	if err != nil {
		panic(err)
	}
}
