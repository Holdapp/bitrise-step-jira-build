package bitrise

type Paging struct {
	Next           string `json:"next"`
	PageItemLimit  int    `json:"page_item_limit"`
	TotalItemCount int    `json:"total_item_count"`
}
