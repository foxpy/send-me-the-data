package flashtest

import (
	"strings"

	"github.com/foxpy/send-me-the-data/cmd/server/flash"
	"golang.org/x/net/html"
)

type Flash struct {
	Kind flash.FlashKind
	Text string
	// TODO: maybe also deal with age
}

func FindAllFlashes(doc *html.Node) (flashes []Flash) {
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode {
			for _, a := range n.Attr {
				if a.Key == "class" {
					for kind := range flash.FlashesCount {
						if a.Val == kind.String() {
							flashes = append(flashes, Flash{
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
