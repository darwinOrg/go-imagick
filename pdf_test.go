package dgimk_test

import (
	dgctx "github.com/darwinOrg/go-common/context"
	dgimk "github.com/darwinOrg/go-imagick"
	dglogger "github.com/darwinOrg/go-logger"
	"testing"
)

func TestConvertPdfToImage(t *testing.T) {
	ctx := &dgctx.DgContext{TraceId: "123"}
	_, err := dgimk.ConvertPdfToImage(ctx, "test.pdf", 800, 1212, 200, 100, ".")
	if err != nil {
		dglogger.Error(ctx, err)
	}
}
