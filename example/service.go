package example

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
)

// Service is service interface
type Service interface {
	SetUser(u *User) error
	SetContact(c *Contact) error
	GetUserViews() ([]*UserView, error)
	GetUserView(id string) (*UserView, error)
	Remove() error
}

// ServiceImpl is service implement
type ServiceImpl struct {
	Conn redis.Conn
}

// SetUser is set user
func (s *ServiceImpl) SetUser(u *User) error {
	if s.Conn == nil {
		return errors.New("not initialized redis conn")
	}

	um, err := s.get()
	if err != nil {
		return err
	}

	um[u.ID] = u

	if err := s.set(um); err != nil {
		return err
	}

	log.Printf("[INFO] set user. id: %s\n", u.ID)

	return nil
}

// SetContact is set contact
func (s *ServiceImpl) SetContact(c *Contact) error {
	if s.Conn == nil {
		return errors.New("not initialized redis conn")
	}

	um, err := s.get()
	if err != nil {
		return err
	}

	if v, ok := um[c.UserID]; ok {
		if v.Contacts == nil {
			v.Contacts = map[string]*Contact{c.ID: c}
		} else {
			v.Contacts[c.ID] = c
		}

		um[c.UserID] = v
	} else {
		return fmt.Errorf("unknown user. id: %s", c.UserID)
	}

	if err := s.set(um); err != nil {
		return err
	}

	log.Printf("[INFO] set contact. id: %s\n", c.ID)

	return nil
}

// GetUserViews is get user views
func (s *ServiceImpl) GetUserViews() ([]*UserView, error) {
	if s.Conn == nil {
		return nil, errors.New("not initialized redis conn")
	}

	um, err := s.get()
	if err != nil {
		return nil, err
	}

	users := make([]*UserView, 0)

	for _, u := range um {
		contacts := make([]*ContactView, 0)

		for _, c := range u.Contacts {
			contacts = append(contacts, &ContactView{
				ID:    c.ID,
				Email: c.Email,
			})
		}

		users = append(users, &UserView{
			ID:       u.ID,
			Name:     u.Name,
			Contacts: contacts,
		})
	}

	log.Printf("[INFO] get user views. length: %d\n", len(users))

	return users, nil
}

func (s *ServiceImpl) GetUserView(id string) (*UserView, error) {
	if s.Conn == nil {
		return nil, errors.New("not initialized redis conn")
	}

	um, err := s.get()
	if err != nil {
		return nil, err
	}

	u, ok := um[id]
	if !ok {
		return nil, fmt.Errorf("not found user. id: %s", id)
	}

	contacts := make([]*ContactView, 0)

	for _, c := range u.Contacts {
		contacts = append(contacts, &ContactView{
			ID:    c.ID,
			Email: c.Email,
		})
	}

	log.Printf("[INFO] get user view. id: %s\n", u.ID)

	return &UserView{
		ID:       u.ID,
		Name:     u.Name,
		Contacts: contacts,
	}, nil
}

// Remove is remove key
func (s *ServiceImpl) Remove() error {
	if s.Conn == nil {
		return errors.New("not initialized redis conn")
	}

	if _, err := s.Conn.Do("DEL", "users"); err != nil {
		return err
	}

	log.Print("[WARN] remove users")

	return nil
}

func (s *ServiceImpl) set(um UserMap) error {
	if s.Conn == nil {
		return errors.New("not initialized redis conn")
	}

	bytes, err := json.Marshal(um)
	if err != nil {
		return err
	}

	if _, err := s.Conn.Do("SET", "users", bytes); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) get() (UserMap, error) {
	if s.Conn == nil {
		return nil, errors.New("not initialized redis conn")
	}

	um := make(UserMap)

	exists, err := redis.Int(s.Conn.Do("EXISTS", "users"))
	if err != nil {
		return nil, err
	}
	if exists == 0 {
		return um, nil
	}

	bytes, err := redis.Bytes(s.Conn.Do("GET", "users"))
	if err != nil {
		return nil, err
	}
	if bytes == nil {
		return um, nil
	}

	if err := json.Unmarshal(bytes, &um); err != nil {
		return nil, err
	}

	return um, nil
}
