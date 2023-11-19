package dgimk

import (
	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
	"gopkg.in/gographics/imagick.v3/imagick"
	"time"
)

// ConvertPdfToImage 转换pdf为图片格式
func ConvertPdfToImage(ctx *dgctx.DgContext, mw *imagick.MagickWand, pdfFile string) (*imagick.MagickWand, error) {
	start := time.Now().UnixMilli()
	if err := mw.ReadImage(pdfFile); err != nil {
		dglogger.Errorf(ctx, "[file: %s] 文件读取失败", pdfFile)
		return nil, err
	}

	var pages = int(mw.GetNumberImages())
	dglogger.Infof(ctx, "[file: %s] 文件页数: %d", pdfFile, pages)
	var mws []*imagick.MagickWand

	for i := 0; i < pages; i++ {
		// This being the page offset
		mw.SetIteratorIndex(i)

		// 压平图像，去掉alpha通道，防止JPG中的alpha变黑, 用在ReadImage之后
		if err := mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_REMOVE); err != nil {
			dglogger.Errorf(ctx, "[file: %s] 压平图像失败：%v", pdfFile, err)
		}

		mw.SetImageFormat("jpg")
		mw.SetImageCompression(imagick.COMPRESSION_JPEG)

		// 如果width > height, 就裁剪成两张
		pWidth := mw.GetImageWidth()
		pHeight := mw.GetImageHeight()

		// 需要裁剪
		if pWidth > pHeight {
			mw.ThumbnailImage(pWidth*2, pHeight)

			tempImage := mw.GetImageFromMagickWand()
			// 由于返回的是指针,需要提前初始化,不然写完左半页tempImage就变了
			rightMw := imagick.NewMagickWandFromImage(tempImage)

			// 左半页
			mw.CropImage(pWidth, pHeight, 0, 0)
			mws = append(mws, mw.GetImage())

			// 右半页
			rightMw.SetImageFormat("jpg")
			rightMw.SetImageCompression(imagick.COMPRESSION_JPEG)
			rightMw.CropImage(pWidth, pHeight, int(pWidth), 0)
			mws = append(mws, rightMw)
		} else {
			mw.ThumbnailImage(pWidth, pHeight)
			mws = append(mws, mw.GetImage())
		}
	}

	defer func() {
		for i := 0; i < len(mws); i++ {
			mws[i].Destroy()
		}
	}()

	mw.Clear()
	for i := 0; i < len(mws); i++ {
		mw.AddImage(mws[i])
	}
	mw.SetFirstIterator()

	// 从上到下追加合并图片
	pmw := mw.AppendImages(true)
	cost := time.Now().UnixMilli() - start
	dglogger.Infof(ctx, "[file: %s] ConvertPdfToImage cost：%d ms", pdfFile, cost)

	return pmw, nil
}
