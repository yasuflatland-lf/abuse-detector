package verify

import (
	"context"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

// TODO : Need to implement timeout
// Scrape links from a url with Chrome headless browser
func Scrape(url string, links *[]string) (bool, error) {
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
		return false, err
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
