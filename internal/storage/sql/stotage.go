package sqlstorage

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/astrviktor/banner-rotation/internal/storage"
	_ "github.com/jackc/pgx/stdlib" //nolint
)

type Storage struct {
	dsn string
	db  *sql.DB
}

func New(dsn string) *Storage {
	return &Storage{dsn, nil}
}

func (s *Storage) Connect() error {
	db, err := sql.Open("pgx", s.dsn)
	if err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		log.Println("Попытка соединения к storage...")
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) CreateSlot(description string) (string, error) {
	id := storage.NewID()

	tx, err := s.db.Begin()
	if err != nil {
		return storage.EmptyID, err
	}

	query := `INSERT INTO banner_rotation.slot
    (id, description)
	VALUES ($1, $2);`

	_, err = tx.Exec(query, id, description)
	if err != nil {
		return storage.EmptyID, err
	}

	err = tx.Commit()
	if err != nil {
		return storage.EmptyID, err
	}

	return id, nil
}

func (s *Storage) CreateBanner(description string) (string, error) {
	id := storage.NewID()

	tx, err := s.db.Begin()
	if err != nil {
		return storage.EmptyID, err
	}

	query := `INSERT INTO banner_rotation.banner
    (id, description)
	VALUES ($1, $2);`

	_, err = tx.Exec(query, id, description)
	if err != nil {
		return storage.EmptyID, err
	}

	query = `INSERT INTO banner_rotation.stat
    (banner_id, segment_id, show_count, click_count)
	SELECT $1, id, 0, 0 FROM banner_rotation.segment;`

	_, err = tx.Exec(query, id)
	if err != nil {
		return storage.EmptyID, err
	}

	err = tx.Commit()
	if err != nil {
		return storage.EmptyID, err
	}

	return id, nil
}

func (s *Storage) CreateSegment(description string) (string, error) {
	id := storage.NewID()

	tx, err := s.db.Begin()
	if err != nil {
		return storage.EmptyID, err
	}

	query := `INSERT INTO banner_rotation.segment
    (id, description)
	VALUES ($1, $2);`

	_, err = tx.Exec(query, id, description)
	if err != nil {
		return storage.EmptyID, err
	}

	query = `INSERT INTO banner_rotation.stat
    (banner_id, segment_id, show_count, click_count)
	SELECT id, $1, 0, 0 FROM banner_rotation.banner;`

	_, err = tx.Exec(query, id)
	if err != nil {
		return storage.EmptyID, err
	}

	err = tx.Commit()
	if err != nil {
		return storage.EmptyID, err
	}

	return id, nil
}

func (s *Storage) CreateRotation(rotation storage.Rotation) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	query := `INSERT INTO banner_rotation.rotation
    (slot_id, banner_id)
	VALUES ($1, $2);`

	_, err = tx.Exec(query, rotation.SlotID, rotation.BannerID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteRotation(rotation storage.Rotation) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	query := `DELETE FROM banner_rotation.rotation WHERE slot_id=$1 AND banner_id = $2;`

	_, err = tx.Exec(query, rotation.SlotID, rotation.BannerID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CreateEvent(slotID, bannerID, segmentID string, action storage.ActionType) error {
	event := storage.Event{
		SlotID:    slotID,
		BannerID:  bannerID,
		SegmentID: segmentID,
		Action:    action,
		Date:      time.Now().UTC(),
	}

	// положить event в kafka

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	var query string
	switch action {
	case storage.Show:
		query = `INSERT INTO banner_rotation.event
    (slot_id, banner_id, segment_id, action, date)
	VALUES ($1, $2, $3, 'show', $4);`
	case storage.Click:
		query = `INSERT INTO banner_rotation.event
    (slot_id, banner_id, segment_id, action, date)
	VALUES ($1, $2, $3, 'click', $4);`
	}

	_, err = tx.Exec(query, event.SlotID, event.BannerID, event.SegmentID, event.Date.Format(time.RFC3339))
	if err != nil {
		return err
	}

	switch action {
	case storage.Show:
		query = `UPDATE banner_rotation.stat
    SET show_count = show_count + 1
    WHERE banner_id = $1 AND segment_id = $2;`
	case storage.Click:
		query = `UPDATE banner_rotation.stat
    SET click_count = click_count + 1
    WHERE banner_id = $1 AND segment_id = $2;`
	}

	_, err = tx.Exec(query, event.BannerID, event.SegmentID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetBannersForSlot(slotID string) ([]string, error) {
	bannersID := make([]string, 0)

	query := `SELECT banner_id
	FROM banner_rotation.rotation 
	WHERE slot_id = $1;`

	rows, err := s.db.Query(query, slotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var bannerID string

		err = rows.Scan(&bannerID)

		if errors.Is(err, sql.ErrNoRows) {
			return bannersID, nil
		}

		if err != nil {
			return bannersID, err
		}

		if rows.Err() != nil {
			return bannersID, err
		}

		bannersID = append(bannersID, bannerID)
	}

	return bannersID, nil
}

func (s *Storage) GetStatForBannerAndSegment(bannerID, segmentID string) (storage.Stat, error) {
	var stat storage.Stat
	stat.BannerID = bannerID
	stat.SegmentID = segmentID

	query := `SELECT show_count, click_count
	FROM banner_rotation.stat
	WHERE banner_id = $1 AND segment_id = $2;`

	rows, err := s.db.Query(query, bannerID, segmentID)
	if err != nil {
		return storage.Stat{}, err
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&stat.ShowCount, &stat.ClickCount)
	if errors.Is(err, sql.ErrNoRows) {
		return storage.Stat{}, errors.New("данные не найдены")
	}

	if err != nil {
		return storage.Stat{}, err
	}

	if rows.Err() != nil {
		return storage.Stat{}, err
	}

	return stat, nil
}
