package models

import "time"

type Locker struct {
	Id       uint64    `json:"id"`
	Name     string    `json:"name"`
	Version  int64     `json:"version"`
	Master   string    `json:"master"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
}
