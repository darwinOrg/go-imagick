package dgimk

import (
	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
	"gopkg.in/gographics/imagick.v3/imagick"
)

func DrawHollowRectOnImage(ctx *dgctx.DgContext, sourceImageFile string, destImageFile string, leftTopX, leftTopY, rightBottomX, rightBottomY float64, strokeWidth float64, color string) error {
	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	cw := imagick.NewPixelWand()
	defer cw.Destroy()
	dw := imagick.NewDrawingWand()
	defer dw.Destroy()

	if err := mw.ReadImage(sourceImageFile); err != nil {
		dglogger.Errorf(ctx, "[file: %s] 文件读取失败", sourceImageFile)
		return err
	}

	if err := dw.PushDrawingWand(); err != nil {
		dglogger.Errorf(ctx, "[file: %s] PushDrawingWand error: %v", sourceImageFile, err)
		return err
	}

	cw.SetColor(color)
	dw.SetStrokeColor(cw)

	cw.SetAlpha(0)
	dw.SetFillColor(cw)

	dw.SetStrokeWidth(strokeWidth)
	dw.SetStrokeAntialias(true)
	dw.Rectangle(leftTopX, leftTopY, rightBottomX, rightBottomY)

	if err := dw.PopDrawingWand(); err != nil {
		dglogger.Errorf(ctx, "[file: %s] PopDrawingWand error: %v", sourceImageFile, err)
		return err
	}

	if err := mw.DrawImage(dw); err != nil {
		dglogger.Errorf(ctx, "[file: %s] DrawImage error: %v", sourceImageFile, err)
		return err
	}

	if err := mw.WriteImage(destImageFile); err != nil {
		dglogger.Errorf(ctx, "[file: %s] WriteImage error: %v", sourceImageFile, err)
		return err
	}

	return nil
}
