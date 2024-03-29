package internalhttp

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/astrviktor/banner-rotation/internal/storage"
)

// GET     /status                                : Проверка статуса сервиса
// POST    /banner                                : Добавляет баннер (description из body), возвращает ID
// POST    /slot                                  : Добавляет слот (description из body), возвращает ID
// POST    /segment                               : Добавляет сегмент (description из body), возвращает ID

// POST    /rotation/{slotID}/{bannerID}          : Добавляет баннер в ротацию в данном слоте.
// DELETE  /rotation/{slotID}/{bannerID}          : Удаляет баннер в ротацию в данном слоте.

// POST    /click/{slotID}/{bannerID}/{segmentID} : Засчитать переход
// Увеличивает счетчик переходов на 1 для указанного баннера в данном слоте в указанной группе.

// POST    /choice/{slotID}/{segmentID}           : Возвращает ID баннера который следует показать в данный момент
// в указанном слоте для указанной соц-дем. группы. Увеличивает число показов баннера в группе.

// GET     /stat/{bannerID}/{segmentID}           : Возвращает статистику по показам и переходам по баннеру для сегмента

type ItemType int

const (
	Banner  ItemType = 1
	Slot    ItemType = 2
	Segment ItemType = 3
)

type Server struct {
	addr    string
	wg      *sync.WaitGroup
	srv     *http.Server
	storage storage.Storage
}

func NewServer(host string, port string, storage storage.Storage) *Server {
	return &Server{
		net.JoinHostPort(host, port),
		&sync.WaitGroup{},
		&http.Server{},
		storage,
	}
}

func (s *Server) Start() {
	if err := s.storage.Connect(); err != nil {
		log.Fatalf("Storage Connect(): %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/status", Logging(s.handleStatus))
	mux.HandleFunc("/banner", Logging(s.handleCreateBanner))
	mux.HandleFunc("/slot", Logging(s.handleCreateSlot))
	mux.HandleFunc("/segment", Logging(s.handleCreateSegment))
	mux.HandleFunc("/rotation/", Logging(s.handleRotation))
	mux.HandleFunc("/click/", Logging(s.handleClick))
	mux.HandleFunc("/choice/", Logging(s.handleChoice))
	mux.HandleFunc("/stat/", Logging(s.handleStat))

	s.srv = &http.Server{
		Addr:    s.addr,
		Handler: mux,
	}

	log.Println("http server starting on address: " + s.addr)

	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe(): %v", err)
		}
		log.Println("http server stopped")
	}()
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown(): %v", err)
	}

	s.storage.Close()

	defer cancel()

	// Wait for ListenAndServe goroutine to close.
	s.wg.Wait()
	log.Println("http server gracefully shutdown")
}
