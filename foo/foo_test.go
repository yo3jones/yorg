package foo

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

var content = `foo
bar
fiz
buz
`

func TestFoo(t *testing.T) {
	var (
		err      error
		filename = "test.jsonl"
		data     []byte
	)
	os.Remove(filename)
	if err = ioutil.WriteFile(filename, []byte(content), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(filename)

	if err = doTheThing(filename); err != nil {
		t.Fatal(err)
	}

	if data, err = ioutil.ReadFile(filename); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("\n%s\n", string(data))
}

func doTheThing(filename string) (err error) {
	var (
		readFile      *os.File
		writeFile     *os.File
		bytesRead     int
		bytesWrittern int
		line          int64
		buffer        = make([]byte, 4)
	)

	if readFile, err = os.OpenFile(filename, os.O_RDONLY, 0666); err != nil {
		return err
	}
	defer readFile.Close()

	if writeFile, err = os.OpenFile(filename, os.O_WRONLY, 0666); err != nil {
		return err
	}
	defer writeFile.Close()

	for {
		if bytesRead, err = readFile.Read(buffer); err != nil && err != io.EOF {
			return err
		}

		if bytesRead < 4 {
			break
		}

		bytesWrittern, err = writeFile.WriteAt(
			[]byte(strings.ToUpper(string(buffer))),
			line*4,
		)
		if err != nil {
			return err
		}

		rev := Reverse(strings.TrimSpace(string(buffer)))

		fmt.Printf(
			"%d %s %s\n",
			bytesWrittern,
			strings.TrimSpace(string(buffer)),
			rev,
		)

		if _, err = writeFile.Seek(0, io.SeekEnd); err != nil {
			return err
		}

		// fmt.Fprintf(writeFile, "%s\n", rev)

		if err == io.EOF {
			break
		}

		line++
	}

	return nil
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
