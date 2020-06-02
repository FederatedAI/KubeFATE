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
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

// Repository is the basic interface of the database CRUD
type Repository interface {
	getCollection() string
	FromBson(m *bson.M) (interface{}, error)
	GetUuid() string
}

// Save the object in the database
func Save(repository Repository) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := ConnectDb()
	if err != nil {
		return "", err
	}
	collection := db.Collection(repository.getCollection())
	_, err = collection.InsertOne(ctx, repository)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return repository.GetUuid(), nil
}

// Find find the objects from the database
func Find(repository Repository) ([]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := ConnectDb()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	collection := db.Collection(repository.getCollection())
	cur, err := collection.Find(ctx, bson.M{}) // find
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cur.Close(ctx)
	var persistents []interface{}
	for cur.Next(ctx) {
		// Decode to bson map
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		// Convert bson.M to struct
		r, err := repository.FromBson(&result)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		persistents = append(persistents, r)
	}
	return persistents, nil
}

// FindByUUID find the object from the database via uuid
func FindByUUID(repository Repository, uuid string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := ConnectDb()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	collection := db.Collection(repository.getCollection())
	filter := bson.M{"uuid": uuid}
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cur.Close(ctx)
	var r interface{}
	for cur.Next(ctx) {
		// Decode to bson map
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		// Convert bson.M to struct
		r, err = repository.FromBson(&result)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}
	return r, nil
}

// FindByUUID find the object from the database via uuid
func FindOneByUUID(repository Repository, uuid string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := ConnectDb()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	collection := db.Collection(repository.getCollection())
	filter := bson.M{"uuid": uuid}
	cur := collection.FindOne(ctx, filter)

	if cur.Err() != nil {
		return nil, cur.Err()
	}

	var r interface{}
	// Decode to bson map
	var result bson.M
	err = cur.Decode(&result)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// Convert bson.M to struct
	r, err = repository.FromBson(&result)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return r, nil
}

// UpdateByUUID Update the object in the database via uuid
func UpdateByUUID(repository Repository, uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := ConnectDb()
	if err != nil {
		log.Println(err)
		return err
	}
	collection := db.Collection(repository.getCollection())
	doc, err := ToDoc(repository)
	if err != nil {
		log.Fatal(err)
		return err
	}

	update := bson.D{
		{Key: "$set", Value: doc},
	}
	filter := bson.M{"uuid": uuid}
	collection.FindOneAndUpdate(ctx, filter, update)

	return nil
}

// ToDoc convert object to bson document
func ToDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}

// ToJson convert object to json string
func ToJson(r interface{}) string {
	b, err := json.Marshal(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return ""
	}
	return string(b)
}

// DeleteByUUID delete object from database via uuid
func DeleteByUUID(repository Repository, uuid string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := ConnectDb()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	collection := db.Collection(repository.getCollection())
	filter := bson.D{{Key: "uuid", Value: uuid}}
	deleteResult, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	return deleteResult.DeletedCount, err
}

// DeleteByUUID delete object from database via uuid
func DeleteOneByUUID(repository Repository, uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := ConnectDb()
	if err != nil {
		log.Println(err)
		return err
	}
	collection := db.Collection(repository.getCollection())
	filter := bson.D{{Key: "uuid", Value: uuid}}
	r, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if r.DeletedCount == 0 {
		return errors.New("this record may not exist(DeletedCount==0)")
	}
	return nil
}

// DeleteByFilter delete object from database via Filter
func DeleteByFilter(repository Repository, filter bson.M) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := ConnectDb()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	collection := db.Collection(repository.getCollection())

	deleteResult, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	return deleteResult.DeletedCount, err
}
// FindByFilter find objects from database via custom filter, such as: findByName, findByStatus
func FindByFilter(repository Repository, filter bson.M) ([]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := ConnectDb()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	collection := db.Collection(repository.getCollection())
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cur.Close(ctx)
	var persistents []interface{}
	for cur.Next(ctx) {
		// Decode to bson map
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		// Convert bson.M to struct
		r, err := repository.FromBson(&result)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		persistents = append(persistents, r)
	}
	return persistents, nil
}

// FindByFilter find objects from database via custom filter, such as: findByName, findByStatus
func FindOneByFilter(repository Repository, filter bson.M) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := ConnectDb()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	collection := db.Collection(repository.getCollection())
	cur := collection.FindOne(ctx, filter)
	if cur.Err() != nil {
		log.Println(err)
		return nil, err
	}
	var persistent interface{}

	// Decode to bson map
	var result bson.M
	err = cur.Decode(&result)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// Convert bson.M to struct
	r, err := repository.FromBson(&result)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	persistent = r

	return persistent, nil
}
