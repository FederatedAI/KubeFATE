/*
* Copyright 2019-2020 VMware, Inc.
* 
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
* http://www.apache.org/licenses/LICENSE-2.0
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
* 
*/
package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
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

func TestFindByUUID(t *testing.T) {
	InitConfigForTest()
	user := &User{}
	results, _ := FindByUUID(user, userJustAddedUuid)
	t.Log(ToJson(results))
}

func Test_encryption(t *testing.T) {
	type args struct {
		plaintext string
		salt      string
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
				salt:      "456",
			},
			want: "9a5f5c39f26968a2b8c359eb03f5b5bcb38123c921cf2e078c31c8563bf52f900a1e3076270579d932c4db02fac68052ce52bae57f513afb24f9eb200b0efa72c74ec829083ba1fdd6830539741de4b7a4e04e6cef0f787664d212dfeca8af52335f71906fba07ed1009518634587e7c5be6fbb0cdd4142edbfccb37119e4831114b5198a4eef6a2874d6bef89f97652d9111e4162e69dc540609e063b660e718511040679f7bff0ab44143d7b85b690eea88309ee57b89c3094b2db55f7226e333d065c712c8576eec3960d424e5b8395d1753d63c60fdbf2200b4a0fe0007073b907affa73e5b57abc1e1a965487618f27e97d1387c431c0c66352e7955f68",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encryption(tt.args.plaintext, tt.args.salt); got != tt.want {
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

func TestUser_IsValid(t *testing.T) {
	InitConfigForTest()
	type fields struct {
		Uuid     string
		Username string
		Password string
		Salt     string
		Email    string
		Status   UserStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			fields: fields{
				Username: "Layne",
				Password: "test",
			},
			want: true,
		},
		{
			name: "",
			fields: fields{
				Username: "Layne",
				Password: "admin",
			},
			want: false,
		},
		{
			name: "",
			fields: fields{
				Username: "admin",
				Password: "admin",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				Uuid:     tt.fields.Uuid,
				Username: tt.fields.Username,
				Password: tt.fields.Password,
				Salt:     tt.fields.Salt,
				Email:    tt.fields.Email,
				Status:   tt.fields.Status,
			}
			if got := user.IsValid(); got != tt.want {
				t.Errorf("User.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
