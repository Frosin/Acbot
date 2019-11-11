package types

import (
	"encoding/json"
	"time"

	pb "acbot/proto/mongo"
)

type User struct {
	ID           string/*primitive.ObjectID*/ `json:"id" bson:"_id"`
	ChatId       int64  `json:"chatId" bson:"chatId"`
	FirstName    string `json:"firstName" bson:"firstName"`
	LastName     string `json:"lastName" bson:"lastName"`
	UserName     string `json:"userName" bson:"userName"`
	Role         string `json:"role" bson:"role"`
	Active       bool   `json:"active" bson:"active"`
	DeactiveTime int64  `json:"deactiveTime" bson:"deactiveTime"`
}

type Activation struct {
	ID        string/*primitive.ObjectID*/ `json:"id" bson:"_id"`
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

// function convert activation type from generated proto file to our standart activation type
// Because protobuf does not support the generation of bson tags is used as an intermediary structure
func GetActivationProto2My(protoAct *pb.Activation) (*Activation, error) {
	//var id primitive.ObjectID
	timestamp, err := time.Parse(time.RFC3339, protoAct.GetTimestamp())
	if err != nil {
		return nil, err
	}
	/*strId := protoAct.GetID()
	if len(strId) == 36 {
		id, err = primitive.ObjectIDFromHex(strId[10 : len(strId)-2])
		if err != nil {
			return nil, err
		}
	}*/
	return &Activation{
		ID:        protoAct.GetID(),
		Timestamp: timestamp,
		User:      protoAct.GetUser(),
		Activator: protoAct.GetActivator(),
		Complete:  protoAct.GetComplete(),
		Retry:     protoAct.GetRetry(),
	}, nil
}

// function convert type.Activation to activation from proto genered go file
// Because protobuf does not support the generation of bson tag is used as an intermediary structure
func GetActivationMy2Proto(activation *Activation) *pb.Activation {
	return &pb.Activation{
		ID:        activation.ID,
		Timestamp: activation.Timestamp.Format(time.RFC3339),
		User:      activation.User,
		Activator: activation.Activator,
		Complete:  activation.Complete,
		Retry:     activation.Retry,
	}
}

func GetUserProto2My(protoUser *pb.User) *User {
	/*var id primitive.ObjectID
	var err error
	strId := protoUser.GetID()
	if len(strId) == 36 {
		id, err = primitive.ObjectIDFromHex(strId[10 : len(strId)-2])
		if err != nil {
			return nil, err
		}
	}*/
	return &User{
		ID:           protoUser.ID,
		ChatId:       protoUser.ChatId,
		FirstName:    protoUser.FirstName,
		LastName:     protoUser.LastName,
		UserName:     protoUser.UserName,
		Role:         protoUser.Role,
		Active:       protoUser.Active,
		DeactiveTime: protoUser.DeactiveTime,
	}
}

func GetUserMy2Proto(user *User) *pb.User {
	return &pb.User{
		ID:           user.ID,
		ChatId:       user.ChatId,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		UserName:     user.UserName,
		Role:         user.Role,
		Active:       user.Active,
		DeactiveTime: user.DeactiveTime,
	}
}
