package internalhttp

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	memorystorage "github.com/astrviktor/banner-rotation/internal/storage/memory"
)

// POST     /banner                                : Добавляет баннер (из body)
// GET      /banner/{IDBanner}                     : Возвращает баннер по IDBanner

// POST     /slot                                  : Добавляет слот (из body)
// GET      /slot/{IDSlot}                         : Возвращает слот по IDSlot

// POST     /segment                               : Добавляет сегмент (из body)
// GET      /segment/{IDSegment}                   : Возвращает сегмент по IDSegment

// POST     /rotation/{IDSlot}/{IDBanner}          : Добавляет баннер в ротацию в данном слоте.
// DELETE   /rotation/{IDSlot}/{IDBanner}          : Удаляет баннер в ротацию в данном слоте.
// GET      /rotations                             : Возвращает все ротации

// POST     /click/{IDSlot}/{IDBanner}/{IDSegment} : Засчитать переход
// Увеличивает счетчик переходов на 1 для указанного баннера в данном слоте в указанной группе.

// POST     /choice/{IDSlot}/{IDSegment}           : Возвращает ID баннера который следует показать в данный момент
// в указанном слоте для указанной соц-дем. группы. Увеличивает число показов баннера в группе.

// GET     /stat/{IDBanner}/{IDSegment}           : Возвращает статистику, сколько по баннеру для сегмента
// показов и переходов

type Server struct {
	addr    string
	wg      *sync.WaitGroup
	srv     *http.Server
	storage *memorystorage.Storage
}

func NewServer(host string, port string, storage *memorystorage.Storage) *Server {
	return &Server{
		net.JoinHostPort(host, port),
		&sync.WaitGroup{},
		&http.Server{},
		storage,
	}
}

func (s *Server) Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("/status", Logging(s.handleStatus))

	mux.HandleFunc("/banner", Logging(s.handleCreateBanner))
	mux.HandleFunc("/banner/", Logging(s.handleGetBanner))

	mux.HandleFunc("/slot", Logging(s.handleCreateSlot))
	mux.HandleFunc("/slot/", Logging(s.handleGetSlot))

	mux.HandleFunc("/segment", Logging(s.handleCreateSegment))
	mux.HandleFunc("/segment/", Logging(s.handleGetSegment))

	mux.HandleFunc("/rotation/", Logging(s.handleRotation))
	mux.HandleFunc("/rotations", Logging(s.handleGetRotations))

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
		log.Fatalf("Shutdown(): %v", err)
	}
	defer cancel()

	// Wait for ListenAndServe goroutine to close.
	s.wg.Wait()
	log.Println("http server gracefully shutdown")
}
