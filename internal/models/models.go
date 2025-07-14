package models

// структура для декодирования полученных URL
type URL struct {
	RespURL string `json:"url"`
}

// структура для передачи ссылок в канал
type ChanURLs struct {
	TaskID string
	URL    string
}

// структура с информацией по статусу конкретного URL
type URLInfo struct {
	Status string `json:"status"`
	URL    string `json:"url"`
}

// структура для кодирования статуса по задаче
type ResponceDownload struct {
	URLsInfo    []URLInfo
	DownloadURL string `json:"download_link,omitempty"`
}
