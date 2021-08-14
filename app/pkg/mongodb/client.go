package eobDB

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, host, port, username, password, database, authSource string) (*mongo.Database, error) {
	var eobDB string
	var anonymous bool
	if username == "" || password == "" {
		anonymous = true
		eobDB = fmt.Sprintf("mongodb://%s:%s", host, port)
	} else {
		eobDB = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	}

	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(eobDB)
	if !anonymous {
		clientOptions.SetAuth(options.Credential{
			AuthSource:  authSource,
			Username:    username,
			Password:    password,
			PasswordSet: true,
		})
	}
	client, err := mongo.Connect(reqCtx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create client to eobDB due to error %w", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create client to eobDB due to error %w", err)
	}

	return client.Database(database), nil
}
