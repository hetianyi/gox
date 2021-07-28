package fontx_test

import (
	"bufio"
	"fmt"
	"github.com/golang/freetype"
	"github.com/hetianyi/gox/img"
	"golang.org/x/image/font"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var (
	dpi      float64 = 72                                             // screen resolution in Dots Per Inch
	fontfile         = "C:\\Users\\hehety\\Desktop\\123\\STKAITI.TTF" // filename of the ttf font
	hinting          = "none"                                         // none | full
	size     float64 = 32                                             // font size in points
	spacing  float64 = 1.5                                            // line spacing (e.g. 2 means double spaced)
	wonb             = true                                           // white text on a black background
)

var text = []string{
	"’你好卡萨丁看书看到卡萨丁哪看得看哈喽Twas brillig, and the slithy toves",
	"Did gyre and gimble in the wabe;",
	"你好卡萨丁看书看到卡萨丁哪看得看哈喽,",
	"你好卡萨丁看书看到卡萨丁哪看得看哈喽.",
	"",
	"“Beware the Jabberwock, my son!",
	"The jaws that bite, the claws that catch!",
	"Beware the Jubjub bird, and shun",
	"The frumious Bandersnatch!”",
	"",
	"He took his vorpal sword in hand:",
	"Long time the manxome foe he sought—",
	"So rested he by the Tumtum tree,",
	"And stood awhile in thought.",
	"",
	"And as in uffish thought he stood,",
	"The Jabberwock, with eyes of flame,",
	"Came whiffling through the tulgey wood,",
	"And burbled as it came!",
	"",
	"One, two! One, two! and through and through",
	"The vorpal blade went snicker-snack!",
	"He left it dead, and with its head",
	"He went galumphing back.",
	"",
	"“And hast thou slain the Jabberwock?",
	"Come to my arms, my beamish boy!",
	"O frabjous day! Callooh! Callay!”",
	"He chortled in his joy.",
	"",
	"’Twas brillig, and the slithy toves",
	"Did gyre and gimble in the wabe;",
	"All mimsy were the borogoves,",
	"And the mome raths outgrabe.",
}

func Test1(t *testing.T) {
	// Read the font data.
	fontBytes, err := ioutil.ReadFile(fontfile)
	if err != nil {
		log.Println(err)
		return
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	// Initialize the context.
	fg := image.Black
	if wonb {
		fg = image.White
	}
	base, _ := img.OpenLocalFile("D:\\图片\\【20P】【多分辨率】触摸的质感\\zhigan.png")
	// 以背景图的尺寸创建新画布
	rgba := image.NewRGBA(image.Rect(0, 0, base.GetSource().Bounds().Max.X, base.GetSource().Bounds().Max.Y))

	// 将背景图画到画布
	draw.Draw(rgba, rgba.Bounds(), base.GetSource(), image.ZP, draw.Src)
	// 将freetype绑定到该画布
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(f)
	c.SetFontSize(size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba) // 将freetype绑定到该画布
	c.SetSrc(fg)   // 设置字体颜色
	switch hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	// 绘制文字
	pt := freetype.Pt(10, 10+int(c.PointToFixed(size)>>6))
	for _, s := range text {
		_, err = c.DrawString(s, pt)
		if err != nil {
			log.Println(err)
			return
		}
		pt.Y += c.PointToFixed(size * 2)
	}

	// 保存
	outFile, err := os.Create("C:\\Users\\hehety\\Desktop\\out.png")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("Wrote out.png OK.")
}
