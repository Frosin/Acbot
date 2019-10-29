package types

import (
	"encoding/json"
	"time"
)

type User struct {
	ChatId       int64  `json:"chatId" bson:"chatId"`
	FirstName    string `json:"firstName" bson:"firstName"`
	LastName     string `json:"lastName" bson:"lastName"`
	UserName     string `json:"userName" bson:"userName"`
	Role         string `json:"role" bson:"role"`
	Active       bool   `json:"active" bson:"active"`
	DeactiveTime int64  `json:"deactiveTime" bson:"deactiveTime"`
}

type Activation struct {
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	User      int64     `json:"user" bson:"user"`
	Activator int64     `json:"activator" bson:"activator"`
	Complete  bool      `json:"complete" bson:"complete"`
	Retry     bool      `json:"retry" bson:"retry"`
}

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

func (a *Activation) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}

func (a *Activation) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, a)
}
