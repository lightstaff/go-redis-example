package example

import (
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

func TestService(t *testing.T) {
	conn, err := redis.DialURL("redis://192.168.99.100")
	if err != nil {
		t.Fatal(err)
	}

	s := &ServiceImpl{
		Conn: conn,
	}

	// remove
	if err := s.Remove(); err != nil {
		t.Fatal(err)
	}

	// set user
	uID := uuid.New().String()
	cID1 := uuid.New().String()
	u := &User{
		ID:   uID,
		Name: "test",
		Contacts: map[string]*Contact{cID1: &Contact{
			ID:     cID1,
			Email:  "first@test.com",
			UserID: uID,
		}},
	}

	if err := s.SetUser(u); err != nil {
		t.Fatal(err)
	}

	// set contact
	cID2 := uuid.New().String()
	c := &Contact{
		ID:     cID2,
		Email:  "second@test.com",
		UserID: uID,
	}

	if err := s.SetContact(c); err != nil {
		t.Fatal(err)
	}

	// get user view
	uv, err := s.GetUserView(uID)
	if err != nil {
		t.Fatal(err)
	}

	if uv.ID != u.ID {
		t.Error("user view id not eqaul input user id")
	}

	t.Logf("get user view. data: %v", uv)

	// get user views
	uvs, err := s.GetUserViews()
	if err != nil {
		t.Fatal(err)
	}

	if len(uvs) != 1 {
		t.Error("user views length is not 1")
	}

	t.Logf("get user views. data: %v", uvs)

	// remove
	if err := s.Remove(); err != nil {
		t.Fatal(err)
	}
}
