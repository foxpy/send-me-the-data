package view

import (
	"fmt"

	"github.com/foxpy/send-me-the-data/cmd/server/template"
)

const paginatorWidth = 3

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
	for i := uint(1); i <= totalPages; i++ {
		if !(isCloseTo(i, 1) || isCloseTo(i, totalPages) || isCloseTo(i, currentPage)) {
			pages = append(pages, template.PaginationView{
				Number:        0,
				Link:          "",
				IsPlaceholder: true,
			})
			if i+paginatorWidth-1 < currentPage {
				i = currentPage - paginatorWidth
			} else if i+paginatorWidth-1 < totalPages {
				i = totalPages - paginatorWidth
			}
			continue
		}

		var link string
		if i != currentPage {
			link = fmt.Sprintf("%s?offset=%d", baseUrl, (i-1)*itemsPerPage)
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

	return diff <= paginatorWidth-1
}
