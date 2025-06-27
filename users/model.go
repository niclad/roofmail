package users

import "github.com/uptrace/bun"

type User struct {
	bun.BaseModel

	ID       int64  `bun:",pk,autoincrement" json:"id,omitempty"`
	Username string `bun:",notnull,unique" json:"username"`
	Password string `bun:",notnull" json:"password,omitempty"`
}
