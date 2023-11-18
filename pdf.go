package dgimk

import (
	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/google/uuid"
	"gopkg.in/gographics/imagick.v3/imagick"
)

// ConvertPdfToImage 转换pdf为图片格式
// @resolution: 扫描精度
// @CompressionQuality: 图片质量: 1~100
func ConvertPdfToImage(ctx *dgctx.DgContext, pdfFilePath string, pageWidth uint, pageHeight uint, resolution float64, compressionQuality uint, outImageDir string) (string, error) {
	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	if err := mw.SetResolution(resolution, resolution); err != nil {
		dglogger.Errorf(ctx, "[file: %s] 扫描精度[%f]设置失败", pdfFilePath, resolution)
	}

	if err := mw.ReadImage(pdfFilePath); err != nil {
		dglogger.Errorf(ctx, "[file: %s] 文件读取失败", pdfFilePath)
		return "", err
	}

	var pages = int(mw.GetNumberImages())
	dglogger.Infof(ctx, "[file: %s] 文件页数: %d", pdfFilePath, pages)
	var mws []*imagick.MagickWand

	for i := 0; i < pages; i++ {
		// This being the page offset
		mw.SetIteratorIndex(i)

		// 压平图像，去掉alpha通道，防止JPG中的alpha变黑,用在ReadImage之后
		if err := mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_REMOVE); err != nil {
			dglogger.Errorf(ctx, "[file: %s] 压平图像失败：%v", pdfFilePath, err)
		}

		mw.SetImageFormat("jpg")
		mw.SetImageCompression(imagick.COMPRESSION_JPEG)
		mw.SetImageCompressionQuality(compressionQuality)

		// 如果width>height ,就裁剪成两张
		pWidth := mw.GetImageWidth()
		pHeight := mw.GetImageHeight()

		// 需要裁剪
		if pWidth > pHeight {
			// mw.ResizeImage(pageWidth*2, pageHeight, imagick.FILTER_UNDEFINED, 1.0)
			mw.ThumbnailImage(pageWidth*2, pageHeight)

			tempImage := mw.GetImageFromMagickWand()
			// 由于返回的是指针,需要提前初始化,不然写完左半页tempImage就变了
			rightMw := imagick.NewMagickWandFromImage(tempImage)

			// 左半页
			mw.CropImage(pageWidth, pageHeight, 0, 0)
			mws = append(mws, mw.GetImage())

			// 右半页
			rightMw.SetImageFormat("jpg")
			rightMw.SetImageCompression(imagick.COMPRESSION_JPEG)
			rightMw.SetImageCompressionQuality(compressionQuality)
			rightMw.CropImage(pageWidth, pageHeight, int(pageWidth), 0)
			mws = append(mws, rightMw)
		} else {
			mw.ThumbnailImage(pageWidth, pageHeight)
			mws = append(mws, mw.GetImage())
		}
	}

	defer func() {
		for i := 0; i < len(mws); i++ {
			mws[i].Destroy()
		}
	}()

	fmw := mws[0]
	for i := 1; i < len(mws); i++ {
		fmw.AddImage(mws[i])
	}
	fmw.SetFirstIterator()
	// 从上到下追加合并图片
	amw := fmw.AppendImages(true)
	defer amw.Destroy()
	outImageFile := outImageDir + "/" + uuid.NewString() + ".jpg"
	if err := amw.WriteImage(outImageFile); err != nil {
		dglogger.Errorf(ctx, "[file: %s] 导出图片文件失败：%v", pdfFilePath, err)
		return "", err
	}

	dglogger.Infof(ctx, "[file: %s] 转换完毕，导出图片文件：%s", pdfFilePath, outImageFile)
	return outImageFile, nil
}
