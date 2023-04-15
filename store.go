package main

import "fmt"

type Store interface {
	GetUsers() ([]User, error)
	GetUserById(int) (*User, error)
	GetUserByName(string) (*User, error)
}

type InMemoryStore struct {
	users []User
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		users: []User{
			{
				Id:   1,
				Name: "Olly",
			},
		},
	}
}

func (ims *InMemoryStore) GetUsers() ([]User, error) {
	return ims.users, nil
}

func (ims *InMemoryStore) GetUserById(id int) (*User, error) {
	for _, user := range ims.users {
		if user.Id == id {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user with id=%d not found", id)
}

func (ims *InMemoryStore) GetUserByName(name string) (*User, error) {
	for _, user := range ims.users {
		if user.Name == name {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user with name=%s not found", name)
}
