package comics

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	host    = "https://m.mh1234.com"
	url     = "https://m.mh1234.com/comic/9329.html"
	imgHost = "https://img.wszwhg.net/"
)

var tokens = make(chan struct{}, 3)

func main() {
	var resp string
	if len(os.Args) > 1 {
		resp = fetch(os.Args[1])
	} else {
		resp = fetch(url)
	}
	if resp == "" {
		os.Exit(1)
	}
	lists := extractList(resp)
	if len(lists) <= 0 {
		fmt.Println("no chapter")
		os.Exit(1)
	}
	wg := sync.WaitGroup{}

	for i := 0; i < len(lists); i++ {
		wg.Add(1)
		go func(i int) {
			tokens <- struct{}{}
			var c string
			for c == "" {
				// fetch until get the info
				c = fetch(lists[i])
				fmt.Printf("fetching %v \n", lists[i])
				time.Sleep(1 * time.Second)
			}
			imgList, title := extractImgList(c)
			if err := os.MkdirAll(title, os.ModePerm); err != nil {
				fmt.Printf("create file error, %s, dirname is %s \n", err, title)
				panic(err)
			}
			for k, v := range imgList {
				err := downloadImage(v, path.Join(title, strconv.Itoa(k)+".jpg"))
				for err != nil {
					fmt.Println(err, title)
					err = downloadImage(v, path.Join(title, strconv.Itoa(k)+".jpg"))
					time.Sleep(1 * time.Second)
				}
			}
			wg.Done()
			<-tokens
		}(i)
	}

	wg.Wait()

}

func downloadImage(url, path string) error {
	fmt.Println("downloading Image", url, path)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errors.New("new Request error")
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.117 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("http get error")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 && resp.StatusCode != 304 {
		return errors.New("Http status code:" + string(resp.StatusCode))
	}
	f, err := os.Create(path)
	if err != nil {
		fmt.Printf("error path is %v\n", path)
		return errors.New("create File error")
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return errors.New("copy failed")
	}

	return nil
}

func fetch(url string) string {
	fmt.Println("Fetch Url", url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("newrequest err:", err)
		return ""
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.117 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Http get err:", err)
		return ""
	}
	if resp.StatusCode != 200 {
		fmt.Println("Http status code:", resp.StatusCode)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read error", err)
		return ""
	}
	return strings.Replace(string(body), "\n", "", -1)
}

func extractList(s string) []string {
	ulRp := regexp.MustCompile(`<ul id="chapter-list-1".*?>(.*?)</ul>`)
	// aRp := regexp.MustCompile(`<a href="(.*?)".*?>`)
	aRp := regexp.MustCompile(`<a href="(.*?)".*?>`)

	ul := ulRp.FindAllStringSubmatch(s, -1)[0][1]
	as := aRp.FindAllStringSubmatch(ul, -1)
	res := make([]string, len(as))
	for k, v := range as {
		res[k] = host + v[1]
	}

	return res

}

func extractImgList(s string) ([]string, string) {
	// get title
	titleRp := regexp.MustCompile(`var pageTitle = "(.*?)-`)

	titleMatched := titleRp.FindAllStringSubmatch(s, -1)
	if len(titleMatched) <= 0 {
		panic("title matched failed")
	}
	title := strings.Replace(titleMatched[0][1], " ", "", -1)

	// get imgpre from chapterPath
	preRp := regexp.MustCompile(`var chapterPath = "(.*?)"`)
	preMatched := preRp.FindAllStringSubmatch(s, -1)
	if len(preMatched) <= 0 {
		panic("chapterPath matched failed")
	}
	pre := preMatched[0][1]

	// get imagelist from scripts
	imgRp := regexp.MustCompile(`var chapterImages = \[(.*?)\]`)
	// fmt.Println(imgRp.FindAllStringSubmatch(s, -1))
	imgMatched := imgRp.FindAllStringSubmatch(s, -1)
	if len(imgMatched) <= 0 {
		panic("imageLists matched failed")
	}
	res := strings.Split(imgMatched[0][1], ",")
	for k, v := range res {
		res[k] = imgHost + pre + strings.Trim(v, "\"")
		res[k] = strings.Replace(res[k], "\\", "", -1)
	}
	title = format(title)

	// sometimes title will return ""
	if title == "" {
		title = "others"
	}
	return res, title
}

func format(s string) string {
	a := []string{"?", "*", "/", "\\", "<", ">", ":", "\"", "|", " ", "."}
	for _, v := range a {
		s = strings.Replace(s, v, "", -1)
	}
	return s
}
