package API

type (
	HackerNewsContent struct {
		By          string `json:"by"`
		Descendants int    `json:"descendants"`
		Id          int    `json:"id"`
		Kids        []int  `json:"kids"`
		Score       int    `json:"score"`
		Text        string `json:"text"`
		Time        int    `json:"time"`
		Title       string `json:"title"`
		Type        string `json:"type"`
	}
)
