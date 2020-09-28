package verify

import (
	"context"
	"fmt"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func Scrape(url string) (string, error) {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var res []*cdp.Node
	allHtml := `//a`

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Nodes(allHtml, &res),
	)

	if err != nil {
		log.Error(err)
	}

	fmt.Println(NodeValues(res))

	return "", nil
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
		fmt.Println(n.Attributes)
		val, ret := FindHref(n.Attributes)
		if true == ret {
			vs = append(vs, val)
		}
	}
	return vs
}
