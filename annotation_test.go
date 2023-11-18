package dgimk_test

import (
	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
	"gopkg.in/gographics/imagick.v3/imagick"
	"testing"
)

func TestDrawAnnotationOnImage(t *testing.T) {
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

	if err := dw.PushDrawingWand(); err != nil {
		dglogger.Errorf(ctx, "[file: %s] PushDrawingWand error: %v", sourceImageFile, err)
		return
	}

	leftTopX, leftTopY, rightBottomX, rightBottomY := float64(28), float64(37), float64(302), float64(72)

	cw.SetColor("red")
	dw.SetStrokeColor(cw)

	cw.SetAlpha(0)
	dw.SetFillColor(cw)

	dw.SetStrokeWidth(1)
	dw.SetStrokeAntialias(true)
	dw.Rectangle(leftTopX, leftTopY, rightBottomX, rightBottomY)

	newX := rightBottomX + 10
	newY := leftTopY - 10
	dw.Line(rightBottomX, leftTopY, newX, newY)
	dw.Line(newX, newY, newX, newY+5)
	dw.Line(newX, newY, newX-5, newY)
	dw.SetFont("chinese_cht.ttf")
	dw.SetFontSize(16)
	dw.Annotation(newX+2, newY+5, "批注")

	if err := dw.PopDrawingWand(); err != nil {
		dglogger.Errorf(ctx, "[file: %s] PopDrawingWand error: %v", sourceImageFile, err)
		return
	}

	if err := mw.DrawImage(dw); err != nil {
		dglogger.Errorf(ctx, "[file: %s] DrawImage error: %v", sourceImageFile, err)
		return
	}

	if err := mw.WriteImage("output.jpg"); err != nil {
		dglogger.Errorf(ctx, "[file: %s] WriteImage error: %v", sourceImageFile, err)
		return
	}
}
