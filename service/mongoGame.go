package service

import (
	"context"
	"goworkshop2/game"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoGame struct {
	Resource *mongo.Database
}

func initContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	return ctx, cancel
}

func (state *MongoGame) Join(name string) (game.Character, error) {
	ctx, _ := initContext()
	collection := state.Resource.Collection("Players")
	character := game.NewCharacter(name)
	_, err := collection.InsertOne(ctx, character)

	return character, err
}

func (state *MongoGame) List() ([]game.Character, error) {
	var players []game.Character
	ctx, _ := initContext()
	collection := state.Resource.Collection("Players")
	cursor, err := collection.Find(ctx, bson.M{})

	for cursor.Next(ctx) {
		var user game.Character
		err = cursor.Decode(&user)
		if err != nil {
			logrus.Print(err)
		}
		players = append(players, user)
	}

	if err != nil {
		return nil, err
	}

	return players, nil
}
