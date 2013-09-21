package zip

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

func ExampleArchive() {
	buffer := new(bytes.Buffer)

	wg := new(sync.WaitGroup)
	wg.Add(1)
	progress := make(chan string)
	go func() {
		for p := range progress {
			fmt.Println(p)
		}
		wg.Done()
	}()

	err := Archive("testdata/foo", true, buffer, progress)
	if err != nil {
		panic(err)
	}
	close(progress)
	wg.Wait()

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

	wg := new(sync.WaitGroup)
	wg.Add(1)
	progress := make(chan string)
	go func() {
		for p := range progress {
			fmt.Println(p)
		}
		wg.Done()
	}()

	err = Unarchive(reader, int64(reader.Len()), tmpDir, progress)
	if err != nil {
		panic(err)
	}
	close(progress)
	wg.Wait()

	// Output:
	// foo/bar
	// foo/baz/aaa
}
