// package gifx.
// ref:
//  https://ezgif.com
//  https://stackoverflow.com/questions/33295023/how-to-split-gif-into-images
package gifx

import (
	"bytes"
	"github.com/disintegration/imaging"
	"github.com/hetianyi/gox/file"
	"github.com/hetianyi/gox/img"
	"image"
	"image/draw"
	"image/gif"
	"io"
)

type Gif struct {
	src *gif.GIF
}

// 获取原gif
func (src *Gif) GetSource() *gif.GIF {
	return src.src
}

// NewGif 构造一个新Gif类
func NewGif(gif *gif.GIF) *Gif {
	return &Gif{gif}
}

// LoadFromLocalFile 从文件加载gif
func LoadFromLocalFile(path string) (*Gif, error) {
	f, err := file.GetFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	g, err := gif.DecodeAll(f)
	if err != nil {
		return nil, err
	}
	return NewGif(g), nil
}

// LoadFromReader 从流加载gif
func LoadFromReader(reader io.Reader) (*Gif, error) {
	g, err := gif.DecodeAll(reader)
	if err != nil {
		return nil, err
	}
	return NewGif(g), nil
}

// AddWaterMark adds a watermark to this image.
//
// 为次图像添加水印。
func (src *Gif) AddWaterMark(watermark *img.Image, anchor imaging.Anchor, paddingX int, paddingY int, opacity float64) (*Gif, error) {
	g := src.src
	//imgWidth, imgHeight := getGifDimensions(g)
	overPaintImage := image.NewRGBA(image.Rect(0, 0, g.Image[0].Bounds().Max.X, g.Image[0].Bounds().Max.Y))
	//draw.Draw(overPaintImage, overPaintImage.Bounds(), g.Image[0], image.ZP, draw.Src)

	watermarkImg := watermark.GetSource()
	var buf bytes.Buffer

	for i, frame := range g.Image {
		buf.Reset()
		// calculate watermark point.
		pot := img.CalculatePt(overPaintImage.Bounds().Size(), watermarkImg.Bounds().Size(), anchor, paddingX, paddingY)
		// render watermark.
		img3 := imaging.Overlay(frame, watermarkImg, pot, opacity)
		// draw it.
		draw.Draw(overPaintImage, frame.Bounds(), img3, image.ZP, draw.Over)
		// convert image.NRGBA to image
		if err := gif.Encode(&buf, overPaintImage, nil); err != nil {
			return src, err
		}
		tmpImg, err := gif.Decode(&buf)
		if err != nil {
			return src, err
		}
		g.Image[i] = tmpImg.(*image.Paletted)
	}
	return src, nil
}

// Generate generates gif using the given images.
//
// 用给定的一系列图片生成一张gif。
func CreateGif(images []*img.Image, delay []int, loopCount int) (*Gif, error) {
	g := &gif.GIF{
		LoopCount: loopCount,
		Delay:     delay,
	}
	var buf bytes.Buffer
	for _, img := range images {
		if err := gif.Encode(&buf, img.GetSource(), nil); err != nil {
			return nil, err
		}
		tmpImg, err := gif.Decode(&buf)
		if err != nil {
			return nil, err
		}
		g.Image = append(g.Image, tmpImg.(*image.Paletted))
	}
	return &Gif{g}, nil
}

// ref: https://stackoverflow.com/questions/33295023/how-to-split-gif-into-images
func getGifDimensions(gif *gif.GIF) (x, y int) {
	var lowestX int
	var lowestY int
	var highestX int
	var highestY int
	for _, img := range gif.Image {
		if img.Rect.Min.X < lowestX {
			lowestX = img.Rect.Min.X
		}
		if img.Rect.Min.Y < lowestY {
			lowestY = img.Rect.Min.Y
		}
		if img.Rect.Max.X > highestX {
			highestX = img.Rect.Max.X
		}
		if img.Rect.Max.Y > highestY {
			highestY = img.Rect.Max.Y
		}
	}
	return highestX - lowestX, highestY - lowestY
}
