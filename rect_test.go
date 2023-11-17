package dgimk_test

import (
	dgctx "github.com/darwinOrg/go-common/context"
	dgimk "github.com/darwinOrg/go-imagick"
	"os"
	"testing"
)

func TestDrawRectOnImage(t *testing.T) {
	ctx := &dgctx.DgContext{TraceId: "123"}
	err := dgimk.DrawRectOnImage(ctx, os.Getenv("imageFile"), "output.jpg",
		28, 37, 302, 72,
		"red", 1)
	if err != nil {
		return
	}
}
