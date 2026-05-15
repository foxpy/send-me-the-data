package view

import (
	"container/heap"
	"fmt"

	"github.com/foxpy/send-me-the-data/cmd/server/idb"
	"github.com/foxpy/send-me-the-data/cmd/server/ifs"
	"github.com/foxpy/send-me-the-data/cmd/server/template"
)

func Files(link idb.Link, files []ifs.File, offset, limit uint) []template.FileView {
	files = sortedRange(files, offset, limit)
	userDownloadable := link.UserDownloadable()
	linkID := link.ID()

	fileViews := make([]template.FileView, 0, len(files))
	for _, file := range files {
		userDownloadLink := ""
		if userDownloadable {
			userDownloadLink = fmt.Sprintf("/%s/%s", linkID, file.Name)
		}

		fileViews = append(fileViews, template.FileView{
			Name:              file.Name,
			UploadedAt:        uint64(file.ModTime.UnixMilli()),
			Size:              bytesToHuman(uint64(file.Size)),
			AdminDownloadLink: fmt.Sprintf("/link/%s/file/%s", linkID, file.Name),
			UserDownloadLink:  userDownloadLink,
			DeleteLink:        fmt.Sprintf("/link/%s/file/%s/delete", linkID, file.Name),
		})
	}

	return fileViews
}

type heapItem struct {
	file     ifs.File
	priority int64
	index    int
}

type minHeap []*heapItem

func (h minHeap) Len() int {
	return len(h)
}

func (h minHeap) Less(i, j int) bool {
	return h[i].priority < h[j].priority
}

func (h minHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *minHeap) Push(x any) {
	n := len(*h)
	item := x.(*heapItem)
	item.index = n
	*h = append(*h, item)
}

func (h *minHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	item.index = -1 // for safety
	*h = old[0 : n-1]
	return item
}

func sortedRange(files []ifs.File, offset, limit uint) []ifs.File {
	if offset > uint(len(files)) {
		return nil
	}

	if limit > uint(len(files))-offset {
		limit = uint(len(files)) - offset
	}

	if limit == 0 {
		return nil
	}

	h := make(minHeap, offset+limit)
	for i := range int(offset + limit) {
		h[i] = &heapItem{
			file:     files[i],
			priority: files[i].ModTime.UnixNano(),
			index:    i,
		}
	}
	heap.Init(&h)
	files = files[offset+limit:]

	for i := range files {
		min := h[0]
		if !(files[i].ModTime.UnixNano() < min.file.ModTime.UnixNano()) {
			heap.Pop(&h)
			heap.Push(&h, &heapItem{
				file:     files[i],
				priority: files[i].ModTime.UnixNano(),
			})
		}
	}

	files = make([]ifs.File, limit)
	for i := range int(limit) {
		files[len(files)-1-i] = heap.Pop(&h).(*heapItem).file
	}

	return files
}
