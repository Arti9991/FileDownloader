package handlers

type HandlersData struct {
	StorageAddr string
	Tasks       []string
}

func InitHandlersData(Stor string) HandlersData {
	HD := new(HandlersData)
	HD.StorageAddr = Stor
	HD.Tasks = make([]string, 0, 3)
	return *HD
}

type URL struct {
	RespURL string `json:"url"`
}
