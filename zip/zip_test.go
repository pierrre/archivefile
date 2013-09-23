package zip

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

func ExampleArchive() {
	buffer := new(bytes.Buffer)

	progress := func(filePath string) {
		fmt.Println(filePath)
	}

	err := Archive("testdata/foo", buffer, progress)
	if err != nil {
		panic(err)
	}

	// Output:
	// foo/bar
	// foo/baz/aaa
}

func ExampleUnarchive() {
	data, err := ioutil.ReadFile("testdata/foo.zip")
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(data)

	tmpDir, err := ioutil.TempDir("", "test_zip")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	progress := func(filePath string) {
		fmt.Println(filePath)
	}

	err = Unarchive(reader, int64(reader.Len()), tmpDir, progress)
	if err != nil {
		panic(err)
	}

	// Output:
	// foo/bar
	// foo/baz/aaa
}
