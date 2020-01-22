package comics

import (
	"fmt"
	"os"
	"testing"
)

func TestFetch(t *testing.T) {
	if fetch("https://www.baidu.com") == "" {
		t.Error(`fetch www.baidu.com failed`)
	}
	if fetch("https://m.mh1234.com") == "" {
		t.Error(`fetch m.mh1234.com failed`)
	}
	if fetch("https://m.mh1234.com/comic/9329.html") == "" {
		t.Error(`fetch m.mh1234.com/comic/9329.html failed`)
	}
	if fetch("https://m.mh1234.com/comic/15250.html") == "" {
		t.Error(`fetch m.mh1234.com/comic/15250.html failed`)
	}
}

func TestFormat(t *testing.T) {
	if format("oh|my?god/") != "ohmygod" {
		t.Error(`format oh|my?god/ failed`)
	}
	if format("...") != "" {
		t.Error(`format ... failed`)
	}
}

func TestExtractList(t *testing.T) {
	if len(extractList(fetch(url))) <= 0 {
		t.Error(`extract list failed`)
	}

	if len(extractList(fetch("https://m.mh1234.com/comic/9329.html"))) <= 0 {
		t.Error(`extract list failed`)
	}

	if len(extractList(fetch("https://m.mh1234.com/comic/15250.html"))) <= 0 {
		t.Error(`extract list failed`)
	}

}

func TestExtractImgListAndDownloadImage(t *testing.T) {
	l := extractList(fetch("https://m.mh1234.com/comic/15250.html"))

	for k, v := range l {
		fmt.Println(k, v)
	}
	s, title := extractImgList(fetch(l[0]))
	if len(s) <= 0 {
		t.Error(`extract image url lists failed`)
	}
	if title == "" {
		t.Error(`extract title failed`)
	}

	downloadImage(s[0], "test.jpg")
	if _, err := os.Stat("test.jpg"); err != nil {
		t.Error(`downloadImage failed`)
	} else {
		os.Remove("test.jpg")
	}

}
