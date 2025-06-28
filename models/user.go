package models

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel

	ID         int64       `bun:",pk,autoincrement" json:"id,omitempty"`
	Username   string      `bun:",notnull,unique" json:"username"`
	Password   string      `bun:",notnull" json:"password,omitempty"`
	CreateTime time.Time   `bun:",nullzero,default:current_timestamp" json:"createTime,omitempty"`
	LastLogin  time.Time   `bun:",nullzero,default:current_timestamp" json:"lastLogin,omitempty"`
	Preference *Preference `bun:"rel:has-one" json:"preference,omitempty"`
	// Likes      []Like      `bun:"rel:has-many" json:"likes,omitempty"`
}
