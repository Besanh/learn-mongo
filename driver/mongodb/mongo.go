package mongodb

import (
	"context"
	"fmt"
	"mongo-fundamential/common/log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type (
	IMongoClient interface {
		GetClient() *mongo.Client
		GetCollection() *mongo.Collection
		Close() error
		InsertOne(data any) (result *mongo.InsertOneResult, err error)
		InsertMany(data []any) (result *mongo.InsertManyResult, err error)
		FindOne(projection bson.D, data any) (result *mongo.Cursor, err error)
		UpdateOne(filter, data any) (result *mongo.UpdateResult, err error)
		UpdateMany(filter, data []any) (result *mongo.UpdateResult, err error)
	}

	MongoConfg struct {
		Host     string
		Port     int
		Database string
	}

	mongoClient struct {
		config MongoConfg
		client *mongo.Client
	}
)

func NewMongoClient(config MongoConfg) (IMongoClient, error) {
	mongo := &mongoClient{
		config: config,
	}

	if err := mongo.Connect(); err != nil {
		log.Error(err)
		return nil, err
	}

	return mongo, nil
}

func (m *mongoClient) Connect() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)))
	if err != nil {
		log.Error(err)
		return
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Error(err)
		return
	}

	return
}

func (m *mongoClient) GetClient() *mongo.Client {
	return m.client
}

func (m *mongoClient) GetCollection() *mongo.Collection {
	return m.client.Database(m.config.Database).Collection(m.config.Database)
}

func (m *mongoClient) Close() error {
	if err := m.client.Disconnect(context.TODO()); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (m *mongoClient) InsertOne(data any) (result *mongo.InsertOneResult, err error) {
	result, err = m.GetCollection().InsertOne(context.TODO(), data)
	if err != nil {
		return
	}

	return
}

func (m *mongoClient) InsertMany(data []any) (result *mongo.InsertManyResult, err error) {
	result, err = m.GetCollection().InsertMany(context.TODO(), data)
	if err != nil {
		return
	}

	return
}

// To get inserted documents, use data nil
func (m *mongoClient) FindOne(projection bson.D, data any) (result *mongo.Cursor, err error) {
	if projection != nil {
		result, err = m.GetCollection().Find(context.TODO(), data, options.Find().SetProjection(projection))
	} else {
		result, err = m.GetCollection().Find(context.TODO(), data)
	}
	if err != nil {
		return
	}
	return
}

func (m *mongoClient) UpdateOne(filter, data any) (result *mongo.UpdateResult, err error) {
	result, err = m.GetCollection().UpdateOne(context.TODO(), filter, data)
	if err != nil {
		return
	}
	return
}

func (m *mongoClient) UpdateMany(filter, data []any) (result *mongo.UpdateResult, err error) {
	result, err = m.GetCollection().UpdateMany(context.TODO(), filter, data)
	if err != nil {
		return
	}
	return
}

/**
* Delete at most a single document that match a specified filter even though multiple documents may match the specified filter.
 */
func (m *mongoClient) DeleteOne(filter any) (result *mongo.DeleteResult, err error) {
	result, err = m.GetCollection().DeleteOne(context.TODO(), filter)
	if err != nil {
		return
	}
	return
}

// Delete all documents that match a specified filter.
func (m *mongoClient) DeleteMany(filter any) (result *mongo.DeleteResult, err error) {
	result, err = m.GetCollection().DeleteMany(context.TODO(), filter)
	if err != nil {
		return
	}
	return
}
