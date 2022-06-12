package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

var (
	count     int32
	lastCount int32
	c         <-chan time.Time
)

type window struct {
	// 窗口时间戳
	//Timestamp int64
	// 窗口内请求总数
	ReqCount int32
}

type Windows struct {
	// 当前窗口
	flag int
	// 窗口集
	Windows [10]window
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt32(&count, 1)
	_, _ = w.Write([]byte("Hello World!"))
}

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}
	w := Windows{}
	g, _ := errgroup.WithContext(context.Background())
	g.Go(func() error {
		http.HandleFunc("/hello_world", helloWorld)
		return server.ListenAndServe()
	})
	g.Go(func() error {
		sigChan := make(chan os.Signal)
		defer close(sigChan)
		// 这里以 sigterm 为例，实际中可以根据不同的信号做不同的处理
		signal.Notify(sigChan, syscall.SIGTERM)
		_ = <-sigChan
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		// 需要保证服务器优雅终止，但也不能超过一定时间
		return server.Shutdown(ctx)
	})
	g.Go(func() error {
		c = time.Tick(100 * time.Millisecond)
		for next := range c {
			nextT := next.UnixNano()
			//w.Windows[w.flag].Timestamp = nextT
			nowCount := atomic.LoadInt32(&count)
			reqCount := nowCount - lastCount
			w.Windows[w.flag].ReqCount = reqCount
			lastCount = nowCount
			fmt.Printf("windowId: %v, time: %v, total_count: %v, window count: %v \n", w.flag, nextT, count, w.Windows)
			if w.flag == 9 {
				w.flag = 0
			} else {
				w.flag++
			}
		}
		return nil
	})
	err := g.Wait()
	fmt.Println(err)

	// 通过 wrk 方式提交并发请求，比如：wrk -t10 -c30 -d 2s http://127.0.0.1:8080/hello_world
}
