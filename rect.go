package dgimk

import (
	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
	"gopkg.in/gographics/imagick.v3/imagick"
)

func DrawRectOnImage(ctx *dgctx.DgContext, imageFile string, outImageDir string, leftTopX, leftTopY, rightBottomX, rightBottomY float64, color string, strokeWidth float64) error {
	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	cw := imagick.NewPixelWand()
	defer cw.Destroy()
	dw := imagick.NewDrawingWand()
	defer dw.Destroy()

	if err := mw.ReadImage(imageFile); err != nil {
		dglogger.Errorf(ctx, "[file: %s] 文件读取失败", imageFile)
		return err
	}

	if err := dw.PushDrawingWand(); err != nil {
		dglogger.Errorf(ctx, "[file: %s] PushDrawingWand error: %v", imageFile, err)
		return err
	}

	cw.SetColor(color)
	dw.SetStrokeColor(cw)
	dw.SetStrokeWidth(strokeWidth)
	dw.SetStrokeAntialias(true)
	//dw.Rectangle(leftTopX, leftTopY, rightBottomX, rightBottomY)
	dw.Line(leftTopX, leftTopY, leftTopX, rightBottomY)
	dw.Line(leftTopX, leftTopY, rightBottomX, leftTopY)
	dw.Line(rightBottomX, rightBottomY, leftTopX, rightBottomY)
	dw.Line(rightBottomX, rightBottomY, rightBottomX, leftTopY)

	if err := dw.PopDrawingWand(); err != nil {
		dglogger.Errorf(ctx, "[file: %s] PopDrawingWand error: %v", imageFile, err)
		return err
	}

	if err := mw.DrawImage(dw); err != nil {
		dglogger.Errorf(ctx, "[file: %s] DrawImage error: %v", imageFile, err)
		return err
	}

	if err := mw.WriteImage(outImageDir); err != nil {
		dglogger.Errorf(ctx, "[file: %s] WriteImage error: %v", imageFile, err)
		return err
	}

	return nil
}
