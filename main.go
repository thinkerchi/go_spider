package main

import (
	"crawler/models"
	"fmt"
	//	"io"
	"io/ioutil"
	"net/http"
	//	"os"
	"regexp"
)

var (
	crawled map[string]bool
)

func init() {
	fmt.Println("crawled is initializing...")

	crawled = make(map[string]bool)

	fmt.Println("crawled is initialized...")
}

const N = 800

func main() {
	Crawl("http://www.qq.com")
}

func Crawl(s string) {
	fmt.Println("toCrawl is initializing....")
	toCrawl := make(chan string, 4096)
	fmt.Println("toCrawl is initialized......")

	// crawled = make(map[string]bool)
	toCrawl <- s
	fmt.Printf("append %s to toCrawl\n", s)

	c := make(chan int, N)
	for i := 0; i < N; i++ {
		fmt.Println("start a goroutine %d", i)
		go Download(toCrawl, c)
	}

	for i := 0; i < N; i++ {
		fmt.Println("end a goroutine %d", i)
		<-c
	}
}

func Download(chs chan string, c chan<- int) {
	for url := range chs {

		fmt.Println("Get the url: %s", url)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("http.Get() error:", err)
			continue
		}

		fmt.Println("http.Get() is succussful...")

		//		io.Copy(os.Stdout, resp.Body)

		defer resp.Body.Close()
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Read the Body error: ", err)
			continue
		}

		//		fmt.Printf("Read %d bytes from body\n", cnt)

		//		fmt.Println("Get the content: ", string(content))

		if err != nil {
			fmt.Println("failed to read url: %s", url)
			return
		}

		rgx := "<title>(.*)</title>"
		reg := regexp.MustCompile(rgx)
		title := reg.Find(content)
		models.Add(trim_title(string(title)), url)

		crawled[url] = true

		go func() {
			rgx := "href=\"(.*?)\""
			reg := regexp.MustCompile(rgx)
			urls := reg.FindAll(content, -1)

			for _, l := range urls {
				ref := trim1(string(l))
				//				ref := string(l)
				if _, ok := crawled[ref]; ok {
					continue
				}

				chs <- ref
			}
		}()
	}
	c <- 1
}

func trim_title(title string) string {
	arr := []rune(title)
	if len(arr) < 15 {
		fmt.Println("title is null")
		return ""
	}
	return string(arr[7 : len(arr)-8])
}

func trim1(l string) string {
	arr := []rune(l)
	return string(arr[6 : len(arr)-1])
}

func trim(l string) string {
	arr := []rune(l)
	start := 0
	for i := range arr {
		if i == '"' {
			break
		}
		start++
	}
	end := len(arr) - 1
	for ; end >= 0; end-- {
		if arr[end] == '"' {
			break
		}
	}
	arr = arr[start+1 : end]
	return string(arr)
}

// func RetrieveLinks(content []byte, chs chan string) {
// 	rgx := "href=\"(.*?)\""
// 	reg := regexp.MustCompile(rgx)
// 	urls := reg.FindAll(content, -1)

// 	for _, l := range urls {
// 		if _, ok := crawled[string(l)]; ok {
// 			continue
// 		}

// 		chs <- string(l)
// 	}
// }
