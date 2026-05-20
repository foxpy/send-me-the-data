package flashtest

import (
	"strings"

	"github.com/foxpy/send-me-the-data/src/flash"
	"golang.org/x/net/html"
)

type HTMLFlash struct {
	Kind flash.FlashKind
	Text string
}

func FindAllFlashes(doc *html.Node) (flashes []HTMLFlash) {
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode {
			for _, a := range n.Attr {
				if a.Key == "class" {
					for kind := range flash.FlashesCount {
						if a.Val == kind.String() {
							flashes = append(flashes, HTMLFlash{
								Kind: kind,
								Text: strings.TrimSpace(n.FirstChild.Data),
							})
						}
					}
				}
			}
		}
	}
	return
}
