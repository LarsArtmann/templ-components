package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.ExecPath("/nix/store/wbj6k0ikjzqvlbl9s21ilqshycjv3c02-chromium-152.0.7977.75/bin/chromium"),
			chromedp.Flag("window-size", "1280,900"),
		)...)
	defer cancelAlloc()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	ctx, cancelT := context.WithTimeout(ctx, 60*time.Second)
	defer cancelT()

	step := func(name string, a chromedp.Action) chromedp.Action {
		return chromedp.ActionFunc(func(c context.Context) error {
			t0 := time.Now()
			fmt.Printf("%-14s ... ", name)
			if err := a.Do(c); err != nil {
				fmt.Printf("ERR after %s: %v\n", time.Since(t0).Round(time.Millisecond), err)
				return err
			}
			fmt.Printf("ok %s\n", time.Since(t0).Round(time.Millisecond))
			return nil
		})
	}

	var png []byte
	err := chromedp.Run(ctx,
		step("emulate", chromedp.EmulateViewport(1280, 900)),
		step("navigate", chromedp.Navigate("http://localhost:8901/users")),
		step("waitready", chromedp.WaitReady("body", chromedp.ByQuery)),
		step("sleep", chromedp.Sleep(600*time.Millisecond)),
		step("shot", chromedp.FullScreenshot(&png, 92)),
	)
	if err != nil {
		log.Fatalf("FAIL: %v (png=%d bytes)", err, len(png))
	}
	fmt.Println("PNG bytes:", len(png))
}
