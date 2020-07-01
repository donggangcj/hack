package util

import (
	"context"
	"hack/cmd/test"
	"testing"
)

func TestCompressImageByTinyPNGAPI(t *testing.T) {
	test.MakeTmpDirWithAOverSizeImage(t)
	defer test.RemoveTmDir()

	ctx := context.Background()
	// FIXME: Remove private information
	err := CompressImageByTinyPNGAPI(ctx, "../testdate/tmpdir/oversize.png", "rwzrDC0wQxj2ztC2RCsfRWT17tvV9h63")
	if err != nil {
		t.Fatal(err)
	}
}

