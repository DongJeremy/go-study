package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

var (
	Web   = doSearch("web")
	Image = doSearch("image")
	Video = doSearch("video")
)

var apiMap = map[string]string{
	"web":   "https://www.so.com/s?q=%s",
	"image": "http://image.so.com/i?q=%s",
	"video": "https://video.360kan.com/v?q=%s"}

func doSearch(kind string) Search {
	var api = apiMap[kind]

	return func(query string) (Result, error) {

		api := fmt.Sprintf(api, query)

		req, err := http.NewRequest("GET", api, nil)
		if err != nil {
			log.Fatalf("%v", err)
		}

		client := http.DefaultClient
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("%v", err)
		}

		var body string
		if err == nil {
			b, e := ioutil.ReadAll(resp.Body)
			if e == nil {
				body = string(b)
			}
			defer resp.Body.Close()
		}
		return Result{fmt.Sprintf("%s result for %q\n:%s\n", kind, query, body)}, nil
	}
}

type Result struct {
	str string
}

type Search func(query string) (Result, error)

func main() {
	So := func(ctx context.Context, query string) ([]Result, error) {
		g, ctx := errgroup.WithContext(ctx)
		searches := []Search{Web, Image, Video}
		results := make([]Result, len(searches))
		for i, search := range searches {
			i, search := i, search // 这里是关于闭包的一个坑，详细看这里 https://golang.org/doc/faq#closures_and_goroutines
			g.Go(func() error {
				result, err := search(query)
				if err == nil {
					results[i] = result
				}
				return err
			})
		}
		if err := g.Wait(); err != nil {
			return nil, err
		}
		return results, nil
	}

	results, err := So(context.Background(), "golang")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for _, result := range results {
		fmt.Println(result)
	}
}
