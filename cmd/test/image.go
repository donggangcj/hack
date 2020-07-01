package test

import (
	"io"
	"os"
	"testing"
)

func MakeTmpDirWithAOverSizeImage(t *testing.T) {
	if err := os.Mkdir("../testdate/tmpdir", 0777); err != nil {
		t.Fatal(err)
	}
	create, err := os.Create("../testdate/tmpdir/oversize.png")
	if err != nil {
		t.Fatal(err)
	}
	defer create.Close()
	open, err := os.Open("../testdate/compress/simple.png")
	if err != nil {
		t.Fatal(err)
	}
	defer open.Close()
	_, err = io.Copy(create, open)
	create.Sync()
	if err != nil {
		t.Fatal(err)
	}
}


func RemoveTmDir() {
	os.RemoveAll("../testdate/tmpdir")
}
