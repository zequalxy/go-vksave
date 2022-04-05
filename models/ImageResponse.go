package models

type ImageResponse struct {
	Response struct {
		Items []struct {
			MessageID  int `json:"message_id"`
			FromID     int `json:"from_id"`
			Attachment struct {
				Type  string `json:"type"`
				Photo struct {
					AlbumID   int    `json:"album_id"`
					Date      int    `json:"date"`
					ID        int    `json:"id"`
					OwnerID   int    `json:"owner_id"`
					AccessKey string `json:"access_key"`
					Sizes     []struct {
						Height int    `json:"height"`
						URL    string `json:"url"`
						Type   string `json:"type"`
						Width  int    `json:"width,omitempty"`
					} `json:"sizes"`
					Text    string `json:"text"`
					HasTags bool   `json:"has_tags"`
				} `json:"photo"`
			} `json:"attachment"`
		} `json:"items"`
		NextFrom string `json:"next_from"`
		Profiles []struct {
			ID              int    `json:"id"`
			FirstName       string `json:"first_name"`
			LastName        string `json:"last_name"`
			CanAccessClosed bool   `json:"can_access_closed"`
			IsClosed        bool   `json:"is_closed"`
		} `json:"profiles"`
	} `json:"response"`
}
