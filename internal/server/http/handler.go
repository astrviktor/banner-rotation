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

type ResponseError struct {
	Error string `json:"error"`
}

type Response struct {
	Banner storage.Banner `json:"banner"`
	Error  string         `json:"error"`
}

type ResponseRotations struct {
	Rotations []storage.Rotation `json:"rotations"`
	Error     string             `json:"error"`
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
	_, _ = io.WriteString(w, "OK\n")
}

func (s *Server) handleCreateBanner(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.CreateBanner(w, r)
		return
	}
}

func (s *Server) handleGetBanner(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.GetBanner(w, r)
	}
}

/*
curl --request POST 'http://127.0.0.1:8888/banner' \
--header 'Content-Type: application/json' \
--data-raw '{"id": "0d59d804-bfe9-427f-ab37-cac59a0fbcd3", "description": "123"}'
*/

func (s *Server) CreateBanner(w http.ResponseWriter, r *http.Request) {
	buf := make([]byte, r.ContentLength)
	_, err := r.Body.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		w.WriteHeader(http.StatusInternalServerError)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка при получении данных из запроса %s", err)})
		return
	}

	banner := storage.Banner{}
	if err = json.Unmarshal(buf, &banner); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка при конвертации данных из запроса %s", err)})
		return
	}

	err = s.storage.CreateBanner(banner)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка при создании баннера %s", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
}

// curl --request GET 'http://127.0.0.1:8888/banner/0d59d804-bfe9-427f-ab37-cac59a0fbcd3'

func (s *Server) GetBanner(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	params := strings.Split(path, "/")
	if len(params) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка в формате запроса %s", path)})
		return
	}

	IDBanner := params[2]

	banner, ok, err := s.storage.GetBanner(IDBanner)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка при получении баннера %s", err)})
		return
	}

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		WriteResponse(w, &ResponseError{fmt.Sprintf("баннер с id=%s не найден", IDBanner)})
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteResponse(w, &banner)
}

func (s *Server) handleCreateSlot(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.CreateSlot(w, r)
		return
	}
}

func (s *Server) handleGetSlot(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.GetSlot(w, r)
	}
}

/*
curl --request POST 'http://127.0.0.1:8888/slot' \
--header 'Content-Type: application/json' \
--data-raw '{"id": "0d59d804-bfe9-427f-ab37-cac59a0fbcd3", "description": "123 456"}'
*/

func (s *Server) CreateSlot(w http.ResponseWriter, r *http.Request) {
	buf := make([]byte, r.ContentLength)
	_, err := r.Body.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		w.WriteHeader(http.StatusInternalServerError)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка при получении данных из запроса %s", err)})
		return
	}

	slot := storage.Slot{}
	err = json.Unmarshal(buf, &slot)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка при конвертации данных из запроса %s", err)})
		return
	}

	if err = s.storage.CreateSlot(slot); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка при создании слота %s", err)})
		return
	}
	w.WriteHeader(http.StatusOK)
}

// curl --request GET 'http://127.0.0.1:8888/slot/0d59d804-bfe9-427f-ab37-cac59a0fbcd3'

func (s *Server) GetSlot(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	params := strings.Split(path, "/")
	if len(params) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка в формате запроса %s", path)})
		return
	}

	IDSlot := params[2]

	slot, ok, err := s.storage.GetSlot(IDSlot)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка при получении слота %s", err)})
		return
	}

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		WriteResponse(w, &ResponseError{fmt.Sprintf("слот с id=%s не найден", IDSlot)})
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteResponse(w, &slot)
}

func (s *Server) handleCreateSegment(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.CreateSegment(w, r)
		return
	}
}

func (s *Server) handleGetSegment(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.GetSegment(w, r)
	}
}

/*
curl --request POST 'http://127.0.0.1:8888/segment' \
--header 'Content-Type: application/json' \
--data-raw '{"id": "0d59d804-bfe9-427f-ab37-cac59a0fbcd3", "description": "123 456"}'
*/

func (s *Server) CreateSegment(w http.ResponseWriter, r *http.Request) {
	buf := make([]byte, r.ContentLength)
	_, err := r.Body.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		w.WriteHeader(http.StatusInternalServerError)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка при получении данных из запроса %s", err)})
		return
	}
	segment := storage.Segment{}
	err = json.Unmarshal(buf, &segment)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка при конвертации данных из запроса %s", err)})
		return
	}
	err = s.storage.CreateSegment(segment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка при создании сегмента %s", err)})
		return
	}
	w.WriteHeader(http.StatusOK)
}

// curl --request GET 'http://127.0.0.1:8888/segment/0d59d804-bfe9-427f-ab37-cac59a0fbcd3'

func (s *Server) GetSegment(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	params := strings.Split(path, "/")
	if len(params) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка в формате запроса %s", path)})
		return
	}

	IDSegment := params[2]

	segment, ok, err := s.storage.GetSegment(IDSegment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка при получении сегмента %s", err)})
		return
	}

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		WriteResponse(w, &ResponseError{fmt.Sprintf("сегмент с id=%s не найден", IDSegment)})
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteResponse(w, &segment)
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

func (s *Server) handleGetRotations(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.GetRotations(w, r)
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

// curl --request GET 'http://127.0.0.1:8888/rotations'

func (s *Server) GetRotations(w http.ResponseWriter, r *http.Request) {
	rotations, err := s.storage.GetRotation()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка при получении ротаций %s", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteResponse(w, &ResponseRotations{rotations, ""})
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
		IDSlot:   params[2],
		IDBanner: params[3],
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
		IDSlot:   params[2],
		IDBanner: params[3],
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

	idSlot := params[2]
	idBanner := params[3]
	idSegment := params[4]

	err := s.storage.AddEvent(idSlot, idBanner, idSegment, storage.Click)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка при добавлении клика: %s", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
}

// curl --request POST 'http://127.0.0.1:8888/choice/1/2/3'

func (s *Server) Choice(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	params := strings.Split(path, "/")
	if len(params) != 4 {
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(w, &ResponseError{"ошибка в формате запроса"})
		return
	}

	idSlot := params[2]
	idSegment := params[3]

	idBanner, err := core.GetBanner(s.storage, idSlot, idSegment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка при выборе баннера для показа %s", err)})
		return
	}

	err = s.storage.AddEvent(idSlot, idBanner, idSegment, storage.Show)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteResponse(w, &ResponseError{fmt.Sprintf("ошибка при добавлении показа: %s", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
}
