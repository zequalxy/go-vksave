package models

type ImageResponse struct {
	Response struct {
		Items []struct {
			MessageID  int `json:"message_id"`
			FromID     int `json:"from_id"`
			Attachment struct {
				Type  string `json:"-"`
				Photo struct {
					AlbumID   int     `json:"-"`
					Date      int     `json:"-"`
					ID        int     `json:"-"`
					OwnerID   int     `json:"-"`
					AccessKey string  `json:"-"`
					Sizes     []Sizes `json:"sizes"`
					Text      string  `json:"-"`
					HasTags   bool    `json:"-"`
				} `json:"photo"`
			} `json:"attachment"`
		} `json:"items"`
		NextFrom string `json:"next_from"`
		Profiles []struct {
			ID              int    `json:"-"`
			FirstName       string `json:"-"`
			LastName        string `json:"-"`
			CanAccessClosed bool   `json:"-"`
			IsClosed        bool   `json:"-"`
		} `json:"-"`
	} `json:"response"`
}
type Sizes struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Type   string `json:"type"`
	Width  int    `json:"width,omitempty"`
}

type ByHeight []Sizes

func (b ByHeight) Len() int           { return len(b) }
func (b ByHeight) Less(i, j int) bool { return b[i].Height < b[j].Height }
func (b ByHeight) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
