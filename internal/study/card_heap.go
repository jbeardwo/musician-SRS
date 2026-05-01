package study

type CardHeap []Card

func (h CardHeap) Len() int { return len(h) }
func (h CardHeap) Less(i, j int) bool {
	// 1. Sort by Interval first
	if h[i].Interval != h[j].Interval {
		return h[i].Interval < h[j].Interval
	}
	// 2. If intervals are equal, prioritize cards never reviewed (New Cards)
	if !h[i].LastReviewedAt.Valid && h[j].LastReviewedAt.Valid {
		return true
	}
	if h[i].LastReviewedAt.Valid && !h[j].LastReviewedAt.Valid {
		return false
	}
	// 3. Otherwise, oldest review date first
	return h[i].LastReviewedAt.Time.Before(h[j].LastReviewedAt.Time)
}
func (h CardHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *CardHeap) Push(x any) {
	*h = append(*h, x.(Card))
}

func (h *CardHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}
