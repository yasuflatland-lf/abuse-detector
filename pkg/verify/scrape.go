package verify

import (
	"context"
	"strings"

	"github.com/mafredri/cdp/devtool"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

// Scrape links from a Url with Chrome headless browser
// chromedp uses the external API. For more details, please refer the link below.
// https://docs.browserless.io/docs/go.html#docsNav
func Scrape(ctx context.Context, url string, links *[]string) (bool, error) {
	// Create a new goroutine and send request there.
	// The result goes to errCh channel.
	errCh := make(chan error, 1)

	// Use the DevTools HTTP/JSON API to manage targets (e.g. pages, webworkers).
	devt := devtool.New("http://chromedp:9222")
	pt, err := devt.Get(ctx, devtool.Page)
	if err != nil {
		pt, err = devt.Create(ctx)
		if err != nil {
			errCh <- err
		}
	}

	actxt, cancelActxt := chromedp.NewRemoteAllocator(ctx, pt.WebSocketDebuggerURL)
	defer cancelActxt()

	ctxLocal, _ := chromedp.NewContext(actxt) //

	var res []*cdp.Node
	allHtml := `//a`

	go func() {
		err := chromedp.Run(ctxLocal,
			chromedp.Navigate(url),
			chromedp.Nodes(allHtml, &res),
		)

		errCh <- err
	}()

	select {
	case err := <-errCh:
		if err != nil {
			log.Error(err)
			return false, err
		}

	// Timeout or Cancel comes here.
	case <-ctx.Done():
		<-errCh
		return false, ctx.Err()
	}

	// log.Debug(NodeValues(res))
	*links = NodeValues(res)

	return true, nil
}

func FindHref(attrs []string) (string, bool) {
	for _, c := range attrs {
		if strings.HasPrefix(c, "http") {
			return c, true
		}
	}
	return "", false
}

func NodeValues(nodes []*cdp.Node) []string {
	var vs []string
	for _, n := range nodes {
		val, ret := FindHref(n.Attributes)
		if true == ret {
			vs = append(vs, val)
		}
	}
	return vs
}
