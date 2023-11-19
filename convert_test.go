package dgimk_test

import (
	dgctx "github.com/darwinOrg/go-common/context"
	dgimk "github.com/darwinOrg/go-imagick"
	dglogger "github.com/darwinOrg/go-logger"
	"gopkg.in/gographics/imagick.v3/imagick"
	"testing"
)

func TestConvertPdfToImage(t *testing.T) {
	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	ctx := &dgctx.DgContext{TraceId: "123"}
	pmw, err := dgimk.ConvertPdfToImage(ctx, mw, "1.pdf")
	if err != nil {
		dglogger.Error(ctx, err)
	}
	defer pmw.Destroy()

	if err := pmw.WriteImage("output.jpg"); err != nil {
		dglogger.Errorf(ctx, "导出图片文件失败：%v", err)
	}
}
