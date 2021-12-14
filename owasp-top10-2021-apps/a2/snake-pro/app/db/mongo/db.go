package db

import (
	"fmt"
	"time"

	"github.com/globocom/secDevLabs/owasp-top10-2021-apps/a2/snake-pro/app/config"
	"github.com/globocom/secDevLabs/owasp-top10-2021-apps/a2/snake-pro/app/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Collections names used in MongoDB.
var (
	UserCollection = "users"
)

// DB is the struct that represents mongo session.
type DB struct {
	Session *mgo.Session
}

// mongoConfig is the struct that represents mongo configuration.
type mongoConfig struct {
	Address      string
	DatabaseName string
	UserName     string
	Password     string
}

// Database is the interface's database.
type Database interface {
	Insert(obj interface{}, collection string) error
	Search(query bson.M, selectors []string, collection string, obj interface{}) error
	Update(query bson.M, updateQuery interface{}, collection string) error
	UpdateAll(query, updateQuery bson.M, collection string) error
	Upsert(query bson.M, obj interface{}, collection string) (*mgo.ChangeInfo, error)
	SearchOne(query bson.M, selectors []string, collection string, obj interface{}) error
}

// Connect connects to mongo and returns the session.
func Connect() (*DB, error) {
	mongoConfig := config.APIconfiguration.MongoConf
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{mongoConfig.MongoHost},
		Timeout:  time.Second * 60,
		FailFast: true,
		Database: mongoConfig.MongoDBName,
		Username: mongoConfig.MongoUser,
		Password: mongoConfig.MongoPassword,
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		fmt.Println("Error connecting to Mongo:", err)
		return nil, err
	}
	session.SetSafe(&mgo.Safe{WMode: "majority"})

	if err := session.Ping(); err != nil {
		fmt.Println("Error pinging Mongo after connection:", err)
		return nil, err
	}

	return &DB{Session: session}, nil
}

// autoReconnect checks mongo's connection each second and, if an error is found, reconect to it.
func autoReconnect(session *mgo.Session) {
	var err error
	for {
		err = session.Ping()
		if err != nil {
			fmt.Println("Error pinging Mongo in autoReconnect:", err)
			session.Refresh()
			err = session.Ping()
			if err == nil {
				fmt.Println("Reconnect to MongoDB successful.")
			} else {
				fmt.Println("Reconnect to MongoDB failed:", err)
			}
		}
		time.Sleep(time.Second * 1)
	}
}

// Insert inserts a new document.
func (db *DB) Insert(obj interface{}, collection string) error {
	session := db.Session.Clone()
	c := session.DB("").C(collection)
	defer session.Close()
	return c.Insert(obj)
}

// Update updates a single document.
func (db *DB) Update(query, updateQuery interface{}, collection string) error {
	session := db.Session.Clone()
	c := session.DB("").C(collection)
	defer session.Close()
	err := c.Update(query, updateQuery)
	return err
}

// UpdateAll updates all documents that match the query.
func (db *DB) UpdateAll(query, updateQuery interface{}, collection string) error {
	session := db.Session.Clone()
	c := session.DB("").C(collection)
	defer session.Close()
	_, err := c.UpdateAll(query, updateQuery)
	return err
}

// Search searchs all documents that match the query. If selectors are present, the return will be only the chosen fields.
func (db *DB) Search(query bson.M, selectors []string, collection string, obj interface{}) error {
	session := db.Session.Clone()
	defer session.Close()
	c := session.DB("").C(collection)

	var err error
	if selectors != nil {
		selector := bson.M{}
		for _, v := range selectors {
			selector[v] = 1
		}
		err = c.Find(query).Select(selector).All(obj)
	} else {
		err = c.Find(query).All(obj)
	}
	if err == nil && obj == nil {
		err = mgo.ErrNotFound
	}
	return err
}

// SearchOne searchs for the first element that matchs with the given query.
func (db *DB) SearchOne(query bson.M, selectors []string, collection string, obj interface{}) error {
	session := db.Session.Clone()
	defer session.Close()
	c := session.DB("").C(collection)

	var err error
	if selectors != nil {
		selector := bson.M{}
		for _, v := range selectors {
			selector[v] = 1
		}
		err = c.Find(query).Select(selector).One(obj)
	} else {
		err = c.Find(query).One(obj)
	}
	if err == nil && obj == nil {
		err = mgo.ErrNotFound
	}
	return err
}

// Upsert inserts a document or update it if it already exists.
func (db *DB) Upsert(query bson.M, obj interface{}, collection string) (*mgo.ChangeInfo, error) {
	session := db.Session.Clone()
	c := session.DB("").C(collection)
	defer session.Close()
	return c.Upsert(query, obj)
}

// GetUserData queries MongoDB and returns user's data based on its username.
func GetUserData(mapParams map[string]interface{}) (types.UserData, error) {
	userDataResponse := types.UserData{}
	session, err := Connect()
	if err != nil {
		return userDataResponse, err
	}
	userDataQuery := []bson.M{}
	for k, v := range mapParams {
		userDataQuery = append(userDataQuery, bson.M{k: v})
	}
	userDataFinalQuery := bson.M{"$and": userDataQuery}
	err = session.SearchOne(userDataFinalQuery, nil, UserCollection, &userDataResponse)
	return userDataResponse, err
}

// RegisterUser regisiter into MongoDB a new user and returns an error.
func RegisterUser(userData types.UserData) error {
	session, err := Connect()
	if err != nil {
		return err
	}

	newUserData := bson.M{
		"username": userData.Username,
		"password": userData.Password,
		"userID":   userData.UserID,
	}
	err = session.Insert(newUserData, UserCollection)
	return err

}
