package models

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Article struct {
	Id      string `json:"Id,omitempty" bson:"Id,omitempty"`
	Title   string `json:"Title,omitempty" bson:"Title,omitempty"`
	Desc    string `json:"desc,omitempty" bson:"desc,omitempty"`
	Content string `json:"content,omitempty" bson:"content,omitempty"`
}

type ArticleModel struct {
	Client  *mongo.Client
	Context context.Context
}

func (m ArticleModel) All() ([]*Article, error) {
	collection := m.Client.Database("restgame").Collection("articles")
	var results []*Article
	findOptions := options.Find()
	findOptions.SetLimit(100)

	cur, err := collection.Find(m.Context, bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal("List All Articles: %v", err)
	}
	for cur.Next(m.Context) {
		var element Article
		err := cur.Decode(&element)
		if err != nil {
			log.Fatal("List All Articles: Decoding: %v", err)
		}

		results = append(results, &element)
	}
	if err := cur.Err(); err != nil {
		log.Fatal("List All Articles: Cursor Error: %v", err)
	}
	cur.Close(context.TODO())
	return results, nil
}

func (m ArticleModel) Get(id string) (*Article, error) {
	var result Article
	collection := m.Client.Database("restgame").Collection("articles")
	filter := bson.D{{"Id", id}}

	err := collection.FindOne(m.Context, filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	return &result, err
}

func (m ArticleModel) Update(id string, article *Article) error {
	collection := m.Client.Database("restgame").Collection("articles")
	filter := bson.D{{"Id", id}}
	data := bson.M{
		"$set": article,
	}

	result, err := collection.UpdateOne(m.Context, filter, data)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("update result: %v\n", result)
	return err
}

func (m ArticleModel) Create(article *Article) error {
	log.Println("Create article")
	collection := m.Client.Database("restgame").Collection("articles")
	insertResult, err := collection.InsertOne(m.Context, article)
	if err != nil {
		log.Fatal("Create Article: %v", err)
	} else {
		log.Printf("Created document: %v\n", insertResult.InsertedID)
	}
	return nil
}

func (m ArticleModel) Delete(id string) error {
	collection := m.Client.Database("restgame").Collection("articles")
	filter := bson.M{"Id": id}
	res, err := collection.DeleteOne(m.Context, filter)
	if err != nil {
		log.Fatal("Article Delete: %v", err)
	}
	if res.DeletedCount == 0 {
		log.Println("Document not found")
	}
	return nil
}
