package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func main() {
	var g errgroup.Group
	var urls = []string{
		"http://localhost:1234/sleep?sec=1",
		"http://localhost:1234/sleep?sec=2",
		"http://localhost:1234/sleep?sec=3",
		"https://www23.s232323o.com/",
		"http://localhost:1234/sleep?sec=4",
		"http://localhost:1234/sleep?sec=5",
		"http://localhost:1234/sleep?sec=6",
	}
	for _, url := range urls {
		// 用g.Go 开启 goroutine 并行抓取URL
		url := url
		g.Go(func() error {
			// 抓取url
			resp, err := http.Get(url)
			if err == nil {
				fmt.Println(url)
				b, e := ioutil.ReadAll(resp.Body)
				fmt.Println(string(b), e)
				defer resp.Body.Close()
			}
			return err
		})
	}
	// 等待全部的URL抓取完毕
	if err := g.Wait(); err == nil {
		fmt.Println("全部抓取成功.")
	} else {
		fmt.Println("抓取失败，原因是：", err)
	}
}
