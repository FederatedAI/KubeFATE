package db

import (
	"bytes"
	"crypto/md5"
	"fmt"

	"github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	Uuid     string     `json:"uuid"`
	Username string     `json:"username"`
	Password string     `json:"password,omitempty"`
	Email    string     `json:"email"`
	Status   UserStatus `json:"userStatus"`
}

type UserStatus int

const (
	Deprecate_u UserStatus = iota
	Available_u
)

func (s UserStatus) String() string {
	names := []string{
		"Deprecate",
		"Available",
	}

	return names[s]
}

func (s UserStatus) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(s.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func encryption(plaintext string) string {
	secretaries := md5.Sum([]byte(plaintext + "kubefate salt"))
	return fmt.Sprintf("%x", secretaries)
}

func NewUser(username string, password string, email string) *User {
	u := &User{
		Uuid:     uuid.NewV4().String(),
		Username: username,
		Password: encryption(password),
		Email:    email,
		Status:   Deprecate_u,
	}

	return u
}

func (user *User) getCollection() string {
	return "user"
}

func (user *User) GetUuid() string {
	return user.Uuid
}

func (user *User) FromBson(m *bson.M) (interface{}, error) {
	bsonBytes, err := bson.Marshal(m)
	if err != nil {
		return nil, err
	}
	err = bson.Unmarshal(bsonBytes, user)
	if err != nil {
		return nil, err
	}
	return *user, nil
}

func (user *User) IsValid() bool {
	filter := bson.M{"username": user.Username, "password": encryption(user.Password)}
	users, err := FindByFilter(user, filter)
	if err != nil || len(users) == 0 {
		return false
	}
	return true
}

func (user *User) IsExisted() bool {
	filter := bson.M{"username": user.Username}
	users, err := FindByFilter(user, filter)
	if err != nil || len(users) == 0 {
		return false
	}
	return true
}

func (user *User) Update() error {
	user.Password = encryption(user.Password)
	err := UpdateByUUID(user, user.Uuid)
	return err
}
