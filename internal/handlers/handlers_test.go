package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"downloader/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var reqChan = make(chan models.ChanURLs)

var (
	Stor     = "./"
	UserID   = "125"
	HandWG   sync.WaitGroup
	ServWG   sync.WaitGroup
	MapMU    sync.Mutex
	NumTasks = 3
	NumFiles = 3
)

func TestPostTask(t *testing.T) {

	hd := InitHandlersData(Stor, NumTasks, NumFiles, reqChan, &HandWG, &MapMU)

	type want struct {
		answer     string
		statusCode int
	}

	tests := []struct {
		name    string
		request string
		Tasks   map[string]map[string]string
		want    want
	}{
		{
			name:    "Simple request for code 201",
			request: "/task",
			Tasks: map[string]map[string]string{
				"Task1": map[string]string{
					"file1": "DONE"},
			},
			want: want{
				statusCode: 201,
				answer:     "Task ID is:",
			},
		},
		{
			name:    "Requset for full task lists",
			request: "/task",
			Tasks: map[string]map[string]string{
				"Task1": map[string]string{
					"file1": "DONE"},
				"Task2": map[string]string{
					"file2": "DONE"},
				"Task3": map[string]string{
					"file3": "DONE"},
			},
			want: want{
				statusCode: 429,
				answer:     "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.request, nil)
			hd.Tasks = test.Tasks

			w := httptest.NewRecorder()
			h := http.HandlerFunc(PostTask(hd))
			h(w, request)

			result := w.Result()
			assert.Equal(t, test.want.statusCode, result.StatusCode)

			userResult, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			strResult := string(userResult)
			bl := strings.Contains(strResult, test.want.answer)
			assert.True(t, bl)
		})
	}
}

func TestPostAddr(t *testing.T) {

	hd := InitHandlersData(Stor, NumTasks, NumFiles, reqChan, &HandWG, &MapMU)

	type want struct {
		answer     string
		statusCode int
	}

	tests := []struct {
		name    string
		request string
		Tasks   map[string]map[string]string
		body    []models.URL
		want    want
	}{
		{
			name:    "Simple request for code 201",
			request: "/task/Task1",
			Tasks: map[string]map[string]string{
				"Task1": map[string]string{},
			},
			body: []models.URL{
				models.URL{RespURL: "1.ru"},
				models.URL{RespURL: "2.ru"},
				models.URL{RespURL: "3.ru"},
			},

			want: want{
				statusCode: 202,
			},
		},
		{
			name:    "Request when task list is full",
			request: "/task/Task1",
			Tasks: map[string]map[string]string{
				"Task1": map[string]string{
					"File1": "Done",
					"File2": "Done",
					"File3": "Done",
				},
			},
			body: []models.URL{
				models.URL{RespURL: "1.ru"},
				models.URL{RespURL: "2.ru"},
				models.URL{RespURL: "3.ru"},
			},
			want: want{
				statusCode: 429,
			},
		},
		{
			name:    "Simple request with wrong taskID",
			request: "/task/Task2",
			Tasks: map[string]map[string]string{
				"Task1": map[string]string{},
			},
			body: []models.URL{
				models.URL{RespURL: "1.ru"},
				models.URL{RespURL: "2.ru"},
				models.URL{RespURL: "3.ru"},
			},

			want: want{
				statusCode: 400,
			},
		},
		{
			name:    "Empty request",
			request: "/task/Task2",
			Tasks: map[string]map[string]string{
				"Task1": map[string]string{},
			},
			body: []models.URL{},
			want: want{
				statusCode: 400,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bt, err := json.Marshal(test.body)
			require.NoError(t, err)

			request := httptest.NewRequest(http.MethodPost, test.request, bytes.NewReader(bt))
			request.Header.Add("Content-Type", "application/json")
			hd.Tasks = test.Tasks

			w := httptest.NewRecorder()
			h := http.HandlerFunc(PostUrl(hd))
			h(w, request)

			result := w.Result()
			err = result.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, test.want.statusCode, result.StatusCode)
		})
	}
}

func TestGetStatus(t *testing.T) {

	hd := InitHandlersData(Stor, NumTasks, NumFiles, reqChan, &HandWG, &MapMU)

	type want struct {
		answer      string
		statusCode  int
		contentType string
		Statuses    string
	}

	tests := []struct {
		name    string
		request string
		Tasks   map[string]map[string]string
		want    want
	}{
		{
			name:    "Simple request for code 200",
			request: "/info/Task1",
			Tasks: map[string]map[string]string{
				"Task1": map[string]string{
					"File1": "Done",
					"File2": "Done",
					"File3": "Done",
				},
			},

			want: want{
				statusCode:  200,
				contentType: "application/json",
				Statuses:    "Done",
			},
		},
		{
			name:    "Reqeust with no such task",
			request: "/info/Task1",
			Tasks: map[string]map[string]string{
				"Task2": map[string]string{},
			},

			want: want{
				statusCode:  204,
				contentType: "",
				Statuses:    "Done",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodGet, test.request, nil)
			hd.Tasks = test.Tasks

			w := httptest.NewRecorder()
			h := http.HandlerFunc(GetStatus(hd))
			h(w, request)

			result := w.Result()
			assert.Equal(t, test.want.statusCode, result.StatusCode)
			assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))

			if result.StatusCode != 200 {
				return
			}
			var URLInfo models.ResponceDownload
			err := json.NewDecoder(result.Body).Decode(&URLInfo)
			require.NoError(t, err)

			for _, Url := range URLInfo.URLsInfo {
				assert.Equal(t, test.Tasks["Task1"][Url.URL], test.want.Statuses)
			}
			err = result.Body.Close()
			require.NoError(t, err)

		})
	}
}
