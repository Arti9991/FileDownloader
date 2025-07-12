package models

type URL struct {
	RespURL string `json:"url"`
}

type ChanURLs struct {
	TaskID string
	URL    string
}
