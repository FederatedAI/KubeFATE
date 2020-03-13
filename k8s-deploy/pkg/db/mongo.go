package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
var cMuLock sync.Mutex

type kubeFATEDatabase struct {
	db            *mongo.Database
	mongoURL      string
	mongoUsername string
	mongoPassword string
	mongoDatabase string
	updateFlag    bool
}

func (kubeFateDB *kubeFATEDatabase) returnMongoInfo() string {
	return "mongodb://" + kubeFateDB.mongoUsername + ":" + kubeFateDB.mongoPassword + "@" + kubeFateDB.mongoURL
}

func (kubeFateDB *kubeFATEDatabase) setUpdateFlag() {
	kubeFateDB.updateFlag = true
}

func (kubeFateDB *kubeFATEDatabase) resetUpdateFlag() {
	kubeFateDB.updateFlag = false
}

func (kubeFateDB *kubeFATEDatabase) isUpdate() bool {
	return kubeFateDB.updateFlag
}

// KubeFateDB
var kubeFateDB *kubeFATEDatabase = nil

// NewKubeFATEDatabase returns a singleton kubeFATEDatabase struct
func newKubeFATEDatabase() error {

	tmpDd := &kubeFATEDatabase{
		mongoURL:      viper.GetString("mongo.url"),
		mongoUsername: viper.GetString("mongo.username"),
		mongoPassword: viper.GetString("mongo.password"),
		mongoDatabase: viper.GetString("mongo.database"),
		updateFlag:    false}

	kubeFateDB = tmpDd

	log.Debug().Msg("Initial Mongo Client with url: " + viper.GetString("mongo.url"))

	return nil
}

// HandleKubeFATEDatabaseUpdate updates the status of the kubeFATEDatabase
func handleKubeFATEDatabaseUpdate() error {

	// Check kubeFateDb instance
	if kubeFateDB == nil {
		return fmt.Errorf("kubeFATEDatabase instance is not found")
	}

	// start a gorutine to handle the dynamic update
	go func() {
		for {
			// Update Mongo url
			if kubeFateDB.mongoURL != viper.GetString("mongo.url") {
				kubeFateDB.mongoURL = viper.GetString("mongo.url")
				kubeFateDB.setUpdateFlag()
			}

			// Update Mongo username
			if kubeFateDB.mongoUsername != viper.GetString("mongo.username") {
				kubeFateDB.mongoUsername = viper.GetString("mongo.username")
				kubeFateDB.setUpdateFlag()
			}

			// Update Mongo password
			if kubeFateDB.mongoPassword != viper.GetString("mongo.password") {
				kubeFateDB.mongoPassword = viper.GetString("mongo.password")
				kubeFateDB.setUpdateFlag()
			}

			// Update Mongo database
			if kubeFateDB.mongoDatabase != viper.GetString("mongo.database") {
				kubeFateDB.mongoDatabase = viper.GetString("mongo.database")
				kubeFateDB.setUpdateFlag()
			}

			time.Sleep(1 * time.Second)
		}
	}()

	return nil
}

func initKubeFATEDatabase() error {
	cMuLock.Lock()
	defer cMuLock.Unlock()

	// double check for safety
	if kubeFateDB == nil {
		newKubeFATEDatabase()
		err := handleKubeFATEDatabaseUpdate()
		if err != nil {
			return err
		}
	}
	return nil
}

// ConnectDb initials the kubeFATEDatabase instance and return an available mongo connection
func ConnectDb() (*mongo.Database, error) {
	if kubeFateDB == nil {
		err := initKubeFATEDatabase()
		if err != nil {
			return nil, err
		}
	}

	// check the dynamic update flag
	if kubeFateDB.db == nil || kubeFateDB.isUpdate() {
		// need lock to prevent re-entrancy
		cMuLock.Lock()
		defer cMuLock.Unlock()

		// double check for safety
		if kubeFateDB.db != nil && !kubeFateDB.isUpdate() {
			return kubeFateDB.db, nil
		}

		opts := options.Client().ApplyURI(kubeFateDB.returnMongoInfo())
		client, err := mongo.Connect(ctx, opts) // client
		if err != nil {
			log.Error().Err(err).Msg("mongodb connection error, ")
			return nil, err
		}

		kubeFateDB.db = client.Database(kubeFateDB.mongoDatabase)
		log.Debug().Msg("Successfully initialized Mongo client")

		// reset update flag
		if kubeFateDB.isUpdate() {
			kubeFateDB.resetUpdateFlag()
		}
	}

	return kubeFateDB.db, nil
}

// Disconnect the DB
func Disconnect() error {
	return nil
}

// Ping DB
func Ping() error {
	return nil
}

