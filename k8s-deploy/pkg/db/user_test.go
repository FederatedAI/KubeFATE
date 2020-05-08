package db

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"
)

var userJustAddedUuid string

func TestAddUser(t *testing.T) {
	InitConfigForTest()
	u := NewUser("Layne", "test", "email@vmware.com")
	userUuid, err := Save(u)
	if err == nil {
		t.Log(userUuid)
		userJustAddedUuid = userUuid
	}
}

func TestFindUsers(t *testing.T) {
	InitConfigForTest()
	user := &User{}
	results, _ := Find(user)
	t.Log(ToJson(results))
}

func TestIsExisted(t *testing.T) {
	InitConfigForTest()
	u := NewUser("Layne", "", "")
	result := u.IsExisted()
	if result {
		t.Log("User Layne is valid.")
	}
}

func TestIsValid(t *testing.T) {
	InitConfigForTest()
	u := NewUser("Layne", "test", "email@vmware.com")
	result := u.IsValid()
	if result {
		t.Log("User Layne exists.")
	}
}

func TestFindByUUID(t *testing.T) {
	InitConfigForTest()
	user := &User{}
	results, _ := FindByUUID(user, userJustAddedUuid)
	t.Log(ToJson(results))
}

func Test_encryption(t *testing.T) {
	type args struct {
		plaintext string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				plaintext: "123",
			},
			want: "6fcd47b86e4d288a322cd198c72c7f12",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encryption(tt.args.plaintext); got != tt.want {
				t.Errorf("encryption() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserDeleteAll(t *testing.T) {
	InitConfigForTest()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := ConnectDb()
	if err != nil {
		log.Error().Err(err).Msg("ConnectDb")
	}
	collection := db.Collection(new(User).getCollection())
	filter := bson.D{}
	r, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("DeleteMany")
	}
	if r.DeletedCount == 0 {
		log.Error().Msg("this record may not exist(DeletedCount==0)")
	}
	fmt.Println(r)
	return
}
