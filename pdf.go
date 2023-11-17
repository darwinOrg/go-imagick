package dgimk

import (
	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/google/uuid"
	"gopkg.in/gographics/imagick.v3/imagick"
	"strconv"
)

// ConvertPdfToImage 转换pdf为图片格式
// @resolution: 扫描精度
// @CompressionQuality: 图片质量: 1~100
func ConvertPdfToImage(ctx *dgctx.DgContext, pdfFilePath string, pageWidth uint, pageHeight uint, resolution float64, compressionQuality uint, outImageDir string) ([]string, error) {
	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	if err := mw.SetResolution(resolution, resolution); err != nil {
		dglogger.Errorf(ctx, "[file: %s] 扫描精度[%f]设置失败", pdfFilePath, resolution)
		return nil, err
	}

	if err := mw.ReadImage(pdfFilePath); err != nil {
		dglogger.Errorf(ctx, "[file: %s] 文件读取失败", pdfFilePath)
		return nil, err
	}

	var pages = int(mw.GetNumberImages())
	dglogger.Infof(ctx, "[file: %s] 文件页数: %d", pdfFilePath, pages)

	// 裁剪会使页数增加
	addPages := 0
	imageFilePrefix := outImageDir + "/" + uuid.NewString() + "-"
	var imageFiles []string

	for i := 0; i < pages; i++ {
		// This being the page offset
		mw.SetIteratorIndex(i)

		// 压平图像，去掉alpha通道，防止JPG中的alpha变黑,用在ReadImage之后
		if err := mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_REMOVE); err != nil {
			dglogger.Debugf(ctx, "[file: %s] 压平图像失败：%v", pdfFilePath, err)
			return nil, err
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
			leftMw := imagick.NewMagickWandFromImage(tempImage)

			// 左半页
			mw.CropImage(pageWidth, pageHeight, 0, 0)
			imageFilePath := imageFilePrefix + strconv.Itoa(i+addPages) + ".jpg"
			mw.WriteImage(imageFilePath)
			imageFiles = append(imageFiles, imageFilePath)

			// 右半页
			leftMw.SetImageFormat("jpg")
			leftMw.SetImageCompression(imagick.COMPRESSION_JPEG)
			leftMw.SetImageCompressionQuality(compressionQuality)
			leftMw.CropImage(pageWidth, pageHeight, int(pageWidth), 0)
			addPages++
			imageFilePath = imageFilePrefix + strconv.Itoa(i+addPages) + ".jpg"
			leftMw.WriteImage(imageFilePath)
			imageFiles = append(imageFiles, imageFilePath)
			leftMw.Destroy()
		} else {
			mw.ThumbnailImage(pageWidth, pageHeight)
			imageFilePath := imageFilePrefix + strconv.Itoa(i+addPages) + ".jpg"
			mw.WriteImage(imageFilePath)
			imageFiles = append(imageFiles, imageFilePath)
		}
	}

	dglogger.Infof(ctx, "[file: %s] 转换完毕，共%d个图片文件", pdfFilePath, len(imageFiles))

	return imageFiles, nil
}
