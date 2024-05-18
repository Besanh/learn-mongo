package mongodb

import (
	"context"
	"fmt"
	"mongo-fundamential/common/log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// mongo -u root -p Anh123###!123 15.235.185.218:27017/besanh
type (
	IMongoClient interface {
		GetClient() *mongo.Client
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
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("%s:%d/?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+2.2.6", m.config.Host, m.config.Port)))
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
