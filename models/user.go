package models

import (
	"time"

	"github.com/aureleoules/gocaml/db"
	"gopkg.in/mgo.v2/bson"
)

// User structure
type User struct {
	ID bson.ObjectId `json:"id" bson:"_id"`

	DiscordID     string `json:"discord_id" bson:"discord_id"`
	Username      string `json:"username" bson:"username"`
	Discriminator string `json:"discriminator" bson:"discriminator"`

	SuccessCount int `json:"success_count" bson:"success_count"`
	ErrorCount   int `json:"error_count" bson:"error_count"`

	LastEvaluation time.Time `json:"last_evaluation" bson:"last_evaluation"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
}

// IncrementSuccess increments SuccessCount by 1
func (u *User) IncrementSuccess() error {
	db.GOCAML.C(db.UsersCollection).UpdateId(u.ID, bson.M{
		"$set": bson.M{
			"success_count":   u.SuccessCount + 1,
			"last_evaluation": time.Now(),
		},
	})
	return nil
}

// IncrementError increments ErrorCount by 1
func (u *User) IncrementError() error {
	db.GOCAML.C(db.UsersCollection).UpdateId(u.ID, bson.M{
		"$set": bson.M{
			"error_count":     u.ErrorCount + 1,
			"last_evaluation": time.Now(),
		},
	})
	return nil
}

// Create inserts user in db
func (u *User) Create() (User, error) {
	u.ID = bson.NewObjectId()
	u.CreatedAt = time.Now()
	err := db.GOCAML.C(db.UsersCollection).Insert(u)
	return *u, err
}

// GetUser returns user by DiscordID
func GetUser(id string) (User, error) {
	var user User
	err := db.GOCAML.C(db.UsersCollection).Find(bson.M{
		"discord_id": id,
	}).One(&user)
	return user, err
}

// GetUsers returns users
func GetUsers() ([]User, error) {
	var users []User
	err := db.GOCAML.C(db.UsersCollection).Find(bson.M{}).All(&users)
	return users, err
}
