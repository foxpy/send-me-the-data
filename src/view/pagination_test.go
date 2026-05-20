package view

import (
	"reflect"
	"testing"

	"github.com/foxpy/send-me-the-data/src/template"
)

func TestPagination(t *testing.T) {
	for _, tc := range []struct {
		name                                  string
		currentItem, totalItems, itemsPerPage uint
		baseUrl                               string
		res                                   []template.PaginationView
	}{
		{
			name:         "zero items",
			currentItem:  0,
			totalItems:   0,
			itemsPerPage: 10,
			baseUrl:      "myapp",
			res:          []template.PaginationView{{Number: 1}},
		},
		{
			name:         "half full page",
			currentItem:  3,
			totalItems:   5,
			itemsPerPage: 10,
			baseUrl:      "myapp",
			res:          []template.PaginationView{{Number: 1}},
		},
		{
			name:         "two pages, first one current",
			currentItem:  3,
			totalItems:   14,
			itemsPerPage: 10,
			baseUrl:      "myapp",
			res: []template.PaginationView{
				{Number: 1},
				{Number: 2, Link: "myapp?offset=10"},
			},
		},
		{
			name:         "two pages, first one current (on the edge)",
			currentItem:  9,
			totalItems:   14,
			itemsPerPage: 10,
			baseUrl:      "myapp",
			res: []template.PaginationView{
				{Number: 1},
				{Number: 2, Link: "myapp?offset=10"},
			},
		},
		{
			name:         "two pages, second one current",
			currentItem:  13,
			totalItems:   14,
			itemsPerPage: 10,
			baseUrl:      "myapp",
			res: []template.PaginationView{
				{Number: 1, Link: "myapp?offset=0"},
				{Number: 2},
			},
		},
		{
			name:         "two pages, second one current (on the edge)",
			currentItem:  10,
			totalItems:   14,
			itemsPerPage: 10,
			baseUrl:      "myapp",
			res: []template.PaginationView{
				{Number: 1, Link: "myapp?offset=0"},
				{Number: 2},
			},
		},
		{
			name:         "seven pages, first one current",
			currentItem:  3,
			totalItems:   65,
			itemsPerPage: 10,
			baseUrl:      "myapp",
			res: []template.PaginationView{
				{Number: 1},
				{Number: 2, Link: "myapp?offset=10"},
				{Number: 3, Link: "myapp?offset=20"},
				{IsPlaceholder: true},
				{Number: 5, Link: "myapp?offset=40"},
				{Number: 6, Link: "myapp?offset=50"},
				{Number: 7, Link: "myapp?offset=60"},
			},
		},
		{
			name:         "seven pages, last one current",
			currentItem:  62,
			totalItems:   65,
			itemsPerPage: 10,
			baseUrl:      "myapp",
			res: []template.PaginationView{
				{Number: 1, Link: "myapp?offset=0"},
				{Number: 2, Link: "myapp?offset=10"},
				{Number: 3, Link: "myapp?offset=20"},
				{IsPlaceholder: true},
				{Number: 5, Link: "myapp?offset=40"},
				{Number: 6, Link: "myapp?offset=50"},
				{Number: 7},
			},
		},
		{
			name:         "ten pages, first one current",
			currentItem:  3,
			totalItems:   95,
			itemsPerPage: 10,
			baseUrl:      "myapp",
			res: []template.PaginationView{
				{Number: 1},
				{Number: 2, Link: "myapp?offset=10"},
				{Number: 3, Link: "myapp?offset=20"},
				{IsPlaceholder: true},
				{Number: 8, Link: "myapp?offset=70"},
				{Number: 9, Link: "myapp?offset=80"},
				{Number: 10, Link: "myapp?offset=90"},
			},
		},
		{
			name:         "ten pages, last one current",
			currentItem:  93,
			totalItems:   95,
			itemsPerPage: 10,
			baseUrl:      "myapp",
			res: []template.PaginationView{
				{Number: 1, Link: "myapp?offset=0"},
				{Number: 2, Link: "myapp?offset=10"},
				{Number: 3, Link: "myapp?offset=20"},
				{IsPlaceholder: true},
				{Number: 8, Link: "myapp?offset=70"},
				{Number: 9, Link: "myapp?offset=80"},
				{Number: 10},
			},
		},
		{
			name:         "twenty pages, eleventh one current",
			currentItem:  106,
			totalItems:   197,
			itemsPerPage: 10,
			baseUrl:      "myapp",
			res: []template.PaginationView{
				{Number: 1, Link: "myapp?offset=0"},
				{Number: 2, Link: "myapp?offset=10"},
				{Number: 3, Link: "myapp?offset=20"},
				{IsPlaceholder: true},
				{Number: 9, Link: "myapp?offset=80"},
				{Number: 10, Link: "myapp?offset=90"},
				{Number: 11},
				{Number: 12, Link: "myapp?offset=110"},
				{Number: 13, Link: "myapp?offset=120"},
				{IsPlaceholder: true},
				{Number: 18, Link: "myapp?offset=170"},
				{Number: 19, Link: "myapp?offset=180"},
				{Number: 20, Link: "myapp?offset=190"},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pages := Pagination(tc.currentItem, tc.totalItems, tc.itemsPerPage, tc.baseUrl)

			if len(pages) != len(tc.res) {
				t.Fatalf("expected %d pages, got %d", len(tc.res), len(pages))
			}

			for i := range pages {
				if !reflect.DeepEqual(pages[i], tc.res[i]) {
					t.Fatalf(`
page at index %d is incorrect:
expected: %+v
got:      %+v`, i, tc.res[i], pages[i])
				}
			}
		})
	}
}
