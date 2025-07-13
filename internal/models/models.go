package models

type URL struct {
	RespURL string `json:"url"`
}

type ChanURLs struct {
	TaskID string
	URL    string
}

type URLInfo struct {
	Status string `json:"status"`
	URL    string `json:"url"`
}

type ResponceDownload struct {
	URLsInfo    []URLInfo
	DownloadURL string `json:"download_link,omitempty"`
}
