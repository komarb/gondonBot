package bot

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	log "github.com/sirupsen/logrus"
	"image"
	"image/draw"
	"image/gif"
	"image/png"
	"math"
	"os"
)

func MakeDemotivator(picBytes *[]byte, picFormat string, bigText string, smallText string) *bytes.Buffer {
	wd, err := os.Getwd()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("Can't get path")
		return nil
	}
	template, err := gg.LoadImage(wd + "/res/images/dem_template.jpg")
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("Can't get template image")
		return nil
	}

	dc := gg.NewContextForImage(template)

	err = dc.LoadFontFace(wd+"/res/fonts/timesS.ttf", 98.0)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("Can't load Times font")
		return nil
	}
	tbw, tbh := dc.MeasureString(bigText)
	dc.SetRGB(1, 1, 1)
	dc.DrawString(bigText, (1280-tbw)/2, 800+tbh)

	err = dc.LoadFontFace(wd+"/res/fonts/arialS.ttf", 40.0)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("Can't load Arial font")
		return nil
	}
	tsw, tsh := dc.MeasureString(smallText)
	dc.SetRGB(1, 1, 1)
	dc.DrawString(smallText, (1280-tsw)/2, 920+tsh)

	buf := new(bytes.Buffer)
	if picFormat == "gif" {
		buf, err = MakeGIF(picBytes, dc)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Warn("Can't get gif image")
			return nil
		}
	} else {
		pic, err := GetResizedPic(picBytes)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Warn("Can't get resized image")
			return nil
		}
		dc.DrawImage(pic, 128, 83)

		err = png.Encode(buf, dc.Image())
	}

	return buf
	//err = dc.SavePNG("out.png")
	//if err != nil {
	//	log.WithFields(log.Fields{"error": err}).Warn("Can't save result image")
	//	return
	//}

}
func MakeGIF(picBytes *[]byte, dc *gg.Context) (*bytes.Buffer, error) {
	gifPic, err := gif.DecodeAll(bytes.NewReader(*picBytes))
	// Draw template image to gif
	gifBoundsSize := gifPic.Image[0].Bounds().Size()
	templateSize := dc.Image().Bounds().Size()

	log.Info("gidBoundsSize: ", gifBoundsSize)

	scaleX := 1024.0 / float64(gifBoundsSize.X)
	scaleY := 696.0 / float64(gifBoundsSize.Y)

	log.Info("scaleX, scaleY: ", scaleX, scaleY)

	drawPointX := int(math.Floor(128.0 / scaleX))
	drawPointY := int(math.Floor(83.0 / scaleY))

	newTemplateSizeX := uint(math.Round(float64(templateSize.X) / scaleX))
	newTemplateSizeY := uint(math.Round(float64(templateSize.Y)/scaleY)) + 10
	resizedTemplate := resize.Resize(newTemplateSizeX, newTemplateSizeY, dc.Image(), resize.Lanczos3)

	log.Info("resizedTemplateSize: ", newTemplateSizeX, newTemplateSizeY)

	frames := make([]*image.Paletted, 0)

	if err != nil {
		return nil, fmt.Errorf("cant get gif pic: %s", err)
	}
	for i := 0; i < len(gifPic.Image); i++ {
		template := image.NewPaletted(resizedTemplate.Bounds(), gifPic.Image[i].Palette)
		draw.Draw(template, resizedTemplate.Bounds(), resizedTemplate, resizedTemplate.Bounds().Min, draw.Src)

		//resizedImage := resize.Resize(1024, 686, gifPic.Image[i], resize.Lanczos3)
		draw.Draw(template, gifPic.Image[i].Bounds().Add(image.Point{X: drawPointX, Y: drawPointY}), gifPic.Image[i], gifPic.Image[i].Bounds().Min, draw.Src)
		frames = append(frames, template)
		log.Info("Done frame ", i, "out of ", len(gifPic.Image))
	}
	buf := new(bytes.Buffer)
	err = gif.EncodeAll(buf, &gif.GIF{Image: frames, Delay: gifPic.Delay, LoopCount: 0})
	if err != nil {
		return nil, fmt.Errorf("cant encode gif: %s", err)
	}
	return buf, nil
}

func GetResizedPic(picBytes *[]byte) (image.Image, error) {
	pic, _, err := image.Decode(bytes.NewReader(*picBytes))
	if err != nil {
		return nil, fmt.Errorf("cant get resized pic: %s", err)
	}
	resizedPic := resize.Resize(1024, 686, pic, resize.Lanczos3)
	return resizedPic, nil
}

func getImage(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

func saveImage(img image.Image, filePath string) {
	outFile, err := os.Create("out.png")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, img)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Wrote out.png OK.")
}
