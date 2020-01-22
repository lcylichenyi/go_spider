package comics

import (
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
}
