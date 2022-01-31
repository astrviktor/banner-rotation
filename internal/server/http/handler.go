package internalhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/astrviktor/banner-rotation/internal/core"
	"github.com/astrviktor/banner-rotation/internal/storage"
)

type Description struct {
	Description string `json:"description"`
}

type ResponseError struct {
	Error string `json:"error"`
}

type ResponseBanner struct {
	Banner storage.Banner `json:"banner"`
}

type ResponseID struct {
	ID string `json:"id"`
}

type ResponseStat struct {
	ShowCount  int `json:"showCount"`
	ClickCount int `json:"clickCount"`
}

func WriteResponse(w http.ResponseWriter, resp interface{}) {
	resBuf, err := json.Marshal(resp)
	if err != nil {
		log.Println(fmt.Sprintf("response marshal error: %s", err))
	}
	_, err = w.Write(resBuf)
	if err != nil {
		log.Println(fmt.Sprintf("response marshal error: %s", err))
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

// handlers

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "OK")
}

func (s *Server) handleCreateBanner(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.CreateItem(Banner, w, r)
		return
	}
}

func (s *Server) handleCreateSlot(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.CreateItem(Slot, w, r)
		return
	}
}

func (s *Server) handleCreateSegment(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.CreateItem(Segment, w, r)
		return
	}
}

/*
curl --request POST 'http://127.0.0.1:8888/banner' \
--header 'Content-Type: application/json' \
--data-raw '{"description": "123"}'
*/

/*
curl --request POST 'http://127.0.0.1:8888/slot' \
--header 'Content-Type: application/json' \
--data-raw '{"description": "123"}'
*/

/*
curl --request POST 'http://127.0.0.1:8888/segment' \
--header 'Content-Type: application/json' \
--data-raw '{"description": "123"}'
*/

func (s *Server) CreateItem(item ItemType, w http.ResponseWriter, r *http.Request) {
	buf := make([]byte, r.ContentLength)
	_, err := r.Body.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		w.WriteHeader(http.StatusInternalServerError)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка при получении данных из запроса %s", err)})
		return
	}

	description := Description{}
	if err = json.Unmarshal(buf, &description); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка при конвертации данных из запроса %s", err)})
		return
	}

	var id string
	switch item {
	case Banner:
		id, err = s.storage.CreateBanner(description.Description)
	case Slot:
		id, err = s.storage.CreateSlot(description.Description)
	case Segment:
		id, err = s.storage.CreateSegment(description.Description)
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка при создании %s", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteResponse(w, &ResponseID{ID: id})
}

func (s *Server) handleRotation(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.CreateRotation(w, r)
		return
	}

	if r.Method == http.MethodDelete {
		s.DeleteRotation(w, r)
		return
	}
}

func (s *Server) handleClick(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.Click(w, r)
	}
}

func (s *Server) handleChoice(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.Choice(w, r)
	}
}

func (s *Server) handleStat(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.Stat(w, r)
	}
}

// curl --request POST 'http://127.0.0.1:8888/rotation/1/2'

func (s *Server) CreateRotation(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	params := strings.Split(path, "/")
	if len(params) != 4 {
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка в формате запроса %s", path)})
		return
	}

	rotation := storage.Rotation{
		SlotID:   params[2],
		BannerID: params[3],
	}

	err := s.storage.CreateRotation(rotation)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка при создании ротации %s", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
}

// curl --request DELETE 'http://127.0.0.1:8888/rotation/1/2'

func (s *Server) DeleteRotation(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	params := strings.Split(path, "/")
	if len(params) != 4 {
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка в формате запроса %s", path)})
		return
	}

	rotation := storage.Rotation{
		SlotID:   params[2],
		BannerID: params[3],
	}

	err := s.storage.DeleteRotation(rotation)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка при удалении ротации %s", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
}

// curl --request POST 'http://127.0.0.1:8888/click/1/2/3'

func (s *Server) Click(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	params := strings.Split(path, "/")
	if len(params) != 5 {
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка в формате запроса %s", path)})
		return
	}

	slotID := params[2]
	bannerID := params[3]
	segmentID := params[4]

	err := s.storage.CreateEvent(slotID, bannerID, segmentID, storage.Click)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка при добавлении клика: %s", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
}

// curl --request POST 'http://127.0.0.1:8888/choice/1/2'

func (s *Server) Choice(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	params := strings.Split(path, "/")
	if len(params) != 4 {
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(w, &ResponseError{"ошибка в формате запроса"})
		return
	}

	slotID := params[2]
	segmentID := params[3]

	bannerID, err := core.GetBanner(s.storage, slotID, segmentID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка при выборе баннера для показа %s", err)})
		return
	}

	err = s.storage.CreateEvent(slotID, bannerID, segmentID, storage.Show)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка при добавлении показа: %s", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteResponse(w, &ResponseID{ID: bannerID})
}

// curl --request GET 'http://127.0.0.1:8888/stat/1/2'

func (s *Server) Stat(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	params := strings.Split(path, "/")
	if len(params) != 4 {
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(w, &ResponseError{"ошибка в формате запроса"})
		return
	}

	bannerID := params[2]
	segmentID := params[3]

	stat, err := s.storage.GetStatForBannerAndSegment(bannerID, segmentID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteResponse(w, &ResponseError{"ошибка и получении статистики"})
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteResponse(w, &ResponseStat{ShowCount: stat.ShowCount, ClickCount: stat.ClickCount})
}
