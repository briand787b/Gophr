package main

import "time"

type Image struct {
	ID 		string
	UserID 		string
	Name 		string
	Location 	string
	Size 		int64
	CreatedAt 	time.Time
	Description 	string
}

type ImageStore interface {
	Save(image *Image) error
	Find(id string) (*Image, error)
	FindAll(offset int) ([]Image, error)
	FindAllByUser(user *User, offset int) ([]Image, error)
}

