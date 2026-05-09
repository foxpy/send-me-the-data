package view

import (
	"fmt"

	"github.com/foxpy/send-me-the-data/cmd/server/template"
)

func Pagination(
	currentItem, totalItems, itemsPerPage uint,
	baseUrl string,
) []template.PaginationView {
	totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage
	if totalPages == 0 {
		totalPages = 1
	}

	currentPage := (currentItem + itemsPerPage) / itemsPerPage
	if currentPage == 0 {
		currentPage = 1
	}

	pages := make([]template.PaginationView, 0, totalPages)
	collapsing := false
	// TODO: this O(N) algorithm could be O(C)
	for j := range totalPages {
		i := j + 1
		toCollapse := true
		if isCloseTo(i, 1) || isCloseTo(i, totalPages) || isCloseTo(i, currentPage) {
			toCollapse = false
			collapsing = false
		}

		if toCollapse {
			if !collapsing {
				pages = append(pages, template.PaginationView{
					Number:        0,
					Link:          "",
					IsPlaceholder: true,
				})
				collapsing = true
			}
			continue
		}

		var link string
		if i != currentPage {
			link = fmt.Sprintf("%s?offset=%d&limit=%d", baseUrl, j*itemsPerPage, itemsPerPage)
		}

		pages = append(pages, template.PaginationView{
			Number:        i,
			Link:          link,
			IsPlaceholder: false,
		})
	}

	return pages
}

func isCloseTo(x, target uint) bool {
	var diff uint
	if x > target {
		diff = x - target
	} else {
		diff = target - x
	}

	return diff <= 2
}
