package API

type (
	HackerNewsContentSubmission struct {
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

	HackerNewsContentComment struct {
		By     string `json:"by"`
		Id     int    `json:"id"`
		Kids   []int  `json:"kids"`
		Parent int    `json:"parent"`
		Text   string `json:"text"`
		Time   int    `json:"time"`
		Type   string `json:"type"`
	}

	HackerNewsProfile struct {
		About     string `json:"about"`
		Created   string `json:"created"`
		Id        string `json:"id"`
		Karma     int    `json:"karma"`
		Submitted []int  `json:"submitted"`
	}
)
