package gifx_test

import (
	"github.com/disintegration/imaging"
	"github.com/hetianyi/gox/file"
	"github.com/hetianyi/gox/img"
	"github.com/hetianyi/gox/img/gifx"
	"github.com/hetianyi/gox/logger"
	"image/gif"
	"testing"
)

func init() {
	logger.Init(nil)
}

// GIF打水印
func TestGIFWaterMark1(t *testing.T) {
	g, err := gifx.LoadFromLocalFile("D:\\tmp\\4.gif")
	if err != nil {
		logger.Fatal(err)
	}
	watermark, err := img.OpenLocalFile("D:\\tmp\\mark.png")
	if err != nil {
		logger.Fatal(err)
	}
	watermark = watermark.Resize(50, 50, imaging.Lanczos)
	g.AddWaterMark(watermark, imaging.BottomRight, 10, 10, 1)

	of, err := file.CreateFile("D:\\tmp\\4_watermark.gif")
	if err != nil {
		logger.Panic(err)
	}
	defer of.Close()

	err = gif.EncodeAll(of, g.GetSource())
	if err != nil {
		logger.Panic(err)
	}
}

func TestGenerate(t *testing.T) {
	images := make([]*img.Image, 4)

	images[0], _ = img.OpenLocalFile("D:\\图片\\500px\\2a736754730d804bc978c9fb6e2fe586.jpg")
	images[1], _ = img.OpenLocalFile("D:\\图片\\500px\\3c736f213ef7a694fbaf4592914ba0cc.jpg")
	images[2], _ = img.OpenLocalFile("D:\\图片\\500px\\7aa9a7d51982d10624096ca9ee9bea7d.jpg")
	images[3], _ = img.OpenLocalFile("D:\\图片\\500px\\212baa369e5791e31bfab54f2ddefdec.jpg")
	g, err := gifx.CreateGif(images, []int{100, 100, 100, 100}, 0)
	if err != nil {
		logger.Panic(err)
	}
	watermark, err := img.OpenLocalFile("D:\\tmp\\mark.png")
	if err != nil {
		logger.Fatal(err)
	}
	watermark = watermark.Resize(50, 50, imaging.Lanczos)
	g.AddWaterMark(watermark, imaging.BottomRight, 10, 10, 1)

	of, err := file.CreateFile("D:\\tmp\\merge.gif")
	if err != nil {
		logger.Panic(err)
	}
	defer of.Close()

	err = gif.EncodeAll(of, g.GetSource())
	if err != nil {
		logger.Panic(err)
	}
	logger.Info("merge success")
}
