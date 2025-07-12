package handlers

import "downloader/internal/models"

type HandlersData struct {
	StorageAddr string
	Tasks       map[string]map[string]string
	reqChan     chan models.ChanURLs
}

func InitHandlersData(Stor string, reqChan chan models.ChanURLs) HandlersData {
	HD := new(HandlersData)
	HD.StorageAddr = Stor
	HD.Tasks = make(map[string]map[string]string, 3)
	HD.reqChan = reqChan
	return *HD
}
