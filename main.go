package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	_ "time"

	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/roaet/restgame/models"
)

type Env struct {
	articles interface {
		All() ([]*models.Article, error)
		Get(id string) (*models.Article, error)
		Create(*models.Article) error
		Delete(string) error
		Update(string, *models.Article) error
	}
}

func handleRequests(env *Env) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/articles", env.articleCreate).Methods("POST")
	router.HandleFunc("/articles", env.articleIndex)
	router.HandleFunc("/articles/{id}", env.articleDelete).Methods("DELETE")
	router.HandleFunc("/articles/{id}", env.articleUpdate).Methods("PUT")
	router.HandleFunc("/articles/{id}", env.articleSelect)

	log.Fatal(http.ListenAndServe(":10000", router))
}

func (env *Env) articleUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article models.Article
	json.Unmarshal(reqBody, &article)

	err := env.articles.Update(key, &article)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(article)
}

func (env *Env) articleSelect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	article, err := env.articles.Get(key)
	if err != nil {
		log.Fatal("Error selecting")
	}
	json.NewEncoder(w).Encode(article)
}

func (env *Env) articleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	err := env.articles.Delete(key)
	if err != nil {
		log.Fatal("Error deleting")
	}
}

func (env *Env) articleIndex(w http.ResponseWriter, r *http.Request) {
	articles, err := env.articles.All()
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(articles)
}

func (env *Env) articleCreate(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article models.Article
	json.Unmarshal(reqBody, &article)

	err := env.articles.Create(&article)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(article)
}

const uri string = "mongodb+srv://roaetAdmin:0uUQBAI5XIuS2MkO@cluster0.lp1vo.mongodb.net/learningdb?retryWrites=true&w=majority"

func main() {
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	ctx := context.TODO()
	err = mongoClient.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer mongoClient.Disconnect(ctx)

	env := &Env{
		articles: models.ArticleModel{Client: mongoClient, Context: ctx},
	}
	handleRequests(env)
}
