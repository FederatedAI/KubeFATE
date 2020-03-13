package db

import (
	"testing"
)

var userJustAddedUuid string

func TestAddUser(t *testing.T) {
	InitConfigForTest()
	u := NewUser("Layne", "test", "email@vmware.com")
	userUuid, error := Save(u)
	if error == nil {
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


