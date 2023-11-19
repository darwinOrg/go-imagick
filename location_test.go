package dgimk

import (
	dgctx "github.com/darwinOrg/go-common/context"
	"github.com/darwinOrg/go-common/utils"
	dglogger "github.com/darwinOrg/go-logger"
	"gopkg.in/gographics/imagick.v3/imagick"
	"strings"
	"testing"
)

func TestWordsLocation(t *testing.T) {
	keyword := "营养"
	words := "纯臻营养护发素"
	keywordIndex := strings.Index(words, keyword)
	preWords := words[0:keywordIndex]

	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	cw := imagick.NewPixelWand()
	defer cw.Destroy()
	dw := imagick.NewDrawingWand()
	defer dw.Destroy()

	ctx := &dgctx.DgContext{TraceId: "123"}
	sourceImageFile := "11.jpg"

	if err := mw.ReadImage(sourceImageFile); err != nil {
		dglogger.Errorf(ctx, "[file: %s] 文件读取失败", sourceImageFile)
		return
	}

	cw.SetColor("red")
	dw.SetStrokeColor(cw)

	cw.SetAlpha(0)
	dw.SetFillColor(cw)

	dw.SetStrokeWidth(1)
	dw.SetStrokeAntialias(true)

	//dw.SetFont("chinese_cht.ttf")
	metrics1 := mw.QueryFontMetrics(dw, preWords)
	metrics2 := mw.QueryFontMetrics(dw, words)
	metrics3 := mw.QueryFontMetrics(dw, keyword)

	dglogger.Infof(ctx, "metrics1: %s", utils.MustConvertBeanToJsonString(metrics1))
	dglogger.Infof(ctx, "metrics2: %s", utils.MustConvertBeanToJsonString(metrics2))
	dglogger.Infof(ctx, "metrics3: %s", utils.MustConvertBeanToJsonString(metrics3))

	leftTopX, leftTopY, rightBottomX, rightBottomY := float64(28), float64(37), float64(302), float64(72)
	dw.Rectangle(leftTopX+(rightBottomX-leftTopX)*(metrics1.TextWidth/metrics2.TextWidth), leftTopY,
		leftTopX+(rightBottomX-leftTopX)*((metrics1.TextWidth+metrics3.TextWidth)/metrics2.TextWidth), rightBottomY)

	mw.SetImageFormat("jpg")
	mw.SetImageCompression(imagick.COMPRESSION_JPEG)
	if err := mw.DrawImage(dw); err != nil {
		return
	}
	mw.WriteImage("location.jpg")
}
