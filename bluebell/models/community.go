package models

import "time"

type Community struct {
	ID   int64  `json:"community_id" db:"community_id" ` //社区id
	Name string `json:"name" db:"community_name"`
}

type CommunityDetail struct {
	ID           int64     `json:"community_id" db:"community_id" ` //社区id
	Name         string    `json:"name" db:"community_name"`
	Introduction string    `json:"introduction,omitempty" db:"introduction"`
	CreateTime   time.Time `json:"create_time" db:"create_time"`
}
