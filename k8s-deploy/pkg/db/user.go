package db

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/pbkdf2"
	"k8s.io/apimachinery/pkg/util/rand"
)

const saltSize = 128

type User struct {
	Uuid     string     `json:"uuid,omitempty"`
	Username string     `json:"username,omitempty"`
	Password string     `json:"password,omitempty"`
	Salt     string     `json:"salt,omitempty"`
	Email    string     `json:"email,omitempty"`
	Status   UserStatus `json:"userStatus,omitempty"`
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

func encryption(plaintext, salt string) string {
	iterations := 100000
	digest := sha256.New
	secretaries := pbkdf2.Key([]byte(plaintext), []byte(salt), iterations, 256, digest)
	return fmt.Sprintf("%x", secretaries)
}

func NewUser(username string, password string, email string) *User {
	salt := rand.String(saltSize)
	u := &User{
		Uuid:     uuid.NewV4().String(),
		Username: username,
		Password: encryption(password, salt),
		Salt:     salt,
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
	filter := bson.M{"username": user.Username}
	gotUser, err := FindOneByFilter(new(User), filter)
	if err != nil || gotUser == nil {
		return false
	}

	salt := gotUser.(User).Salt

	if gotUser.(User).Password != encryption(user.Password, salt) {
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
	salt := rand.String(saltSize)
	user.Password = encryption(user.Password, salt)
	err := UpdateByUUID(user, user.Uuid)
	return err
}
