package storage

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/Corray333/authbot/internal/types"
)

const (
	PHONE_LIFETIME = 60 * 15
)

func (s *Storage) SetPhone(chatID int64, query *types.CodeQuery) error {
	serialized, err := json.Marshal(query)
	if err != nil {
		return err
	}
	query = nil

	if res := s.Redis.Set(strconv.Itoa(int(chatID)), string(serialized), PHONE_LIFETIME*time.Second); res.Err() != nil {
		return res.Err()
	}
	return nil
}

func (s *Storage) GetPhone(chatID int64) (*types.CodeQuery, error) {
	res := s.Redis.Get(strconv.Itoa(int(chatID)))
	if res.Err() != nil {
		return nil, res.Err()
	}
	var query types.CodeQuery
	if err := json.Unmarshal([]byte(res.Val()), &query); err != nil {
		return nil, err
	}
	return &query, nil
}
