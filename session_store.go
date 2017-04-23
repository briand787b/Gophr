package main

import (
	"io/ioutil"
	"os"
	"encoding/json"
	"fmt"
)

type SessionStore interface {
	Find(string) (*Session, error)
	Save(*Session) error
	Delete(*Session) error
}

type FileSessionStore struct {
	filename string
	Sessions map[string]Session
}

var globalSessionStore SessionStore

func NewFileSessionStore(name string) (*FileSessionStore, error) {
	store := &FileSessionStore{
		Sessions: map[string]Session{},
		filename: name,
	}

	contents, err := ioutil.ReadFile(name)
	if err != nil {
		// If it's a matter of the file not existing, it's ok
		if os.IsNotExist(err) {
			return store, nil
		}
		return nil, err
	}
	err = json.Unmarshal(contents, store)
	if err != nil {
		return nil, err
	}
	for _, v := range store.Sessions {
		fmt.Println(v)
	}
	return store, err
}

func (s *FileSessionStore) Find(id string) (*Session, error) {
	session, exists := s.Sessions[id]
	if !exists {
		return nil, nil
	}

	return &session, nil
}

func (store *FileSessionStore) Save(session *Session) error {
	store.Sessions[session.ID] = *session
	contents, err := json.MarshalIndent(store, "", "	")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(store.filename, contents, 0660)
}

func (store *FileSessionStore) Delete(session *Session) error {
	delete(store.Sessions, session.ID)
	contents, err := json.MarshalIndent(store, "", "	")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(store.filename, contents, 0660)
}