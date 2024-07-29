package translate

import (
	"bytes"
	"os"
	"testing"
)

func TestReadMarkdownFile(t *testing.T) {
	_, err := readMarkdownFile("./test.md")
	if err != nil {
		t.Fatal(err)
	}

	_, err = readMarkdownFile("./test-not-exist.md")
	if err == nil {
		t.Fatal("should be error")
	}
}

func TestWriteMarkdownFile(t *testing.T) {
	err := writeMarkdownFile("./test-write.md", []byte("# test\n## test sub1\n test stub"))
	if err != nil {
		t.Fatal(err)
	}
	data, err := readMarkdownFile("./test-write.md")
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(data, []byte("# test\n## test sub1\n test stub")) != 0 {
		t.Fatal("should be equal")
	}
	defer func() {
		os.Remove("./test-write.md")
	}()
}

func TestRenameEnMarkdownFile(t *testing.T) {
	newName := renameEnMarkdownFile("../pkg/test/test.md")

	if newName != "../pkg/test/test_en.md" {
		t.Fatal("should be equal")
	}
}
