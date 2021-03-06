package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var COMMENTS_SNS_TOPIC string = "arn:aws:sns:us-east-1:189400394091:comments"
var LIKES_SNS_TOPIC string = "arn:aws:sns:us-east-1:189400394091:likes"

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()
	initRoutes(mx, formatter)
	n.UseHandler(mx)
	return n
}



// API Routes
func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/like/count/{photo_id}", likeCount(formatter)).Methods("GET")
	mx.HandleFunc("/like/{photo_id}", addLike(formatter)).Methods("POST", "OPTIONS")
	mx.HandleFunc("/like/{photo_id}", removeLike(formatter)).Methods("DELETE")
	mx.HandleFunc("/comment/{photo_id}/{user_id}", commentList(formatter)).Methods("GET")
	mx.HandleFunc("/comment/{photo_id}", addComment(formatter)).Methods("POST", "OPTIONS")
	mx.HandleFunc("/comment/{photo_id}", removeComment(formatter)).Methods("DELETE")
	mx.HandleFunc("/ping", pingHandler(formatter)).Methods("GET")
}

// Helper Functions
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// API GET Number of likes
func likeCount(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		enableCors(&w)
		var m likecount
		_ = json.NewDecoder(req.Body).Decode(&m)
		params := mux.Vars(req)
		var photo_id string = params["photo_id"]
		_ = photo_id

		fmt.Println(photo_id)
		count := readLikes(photo_id)

		formatter.JSON(w, http.StatusOK, struct{ Count int64 }{count})
	}
}

// API to POST a new like
func addLike(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		setupResponse(&w, req)
		if (*req).Method == "OPTIONS" {
			return
		}
		var m like
		_ = json.NewDecoder(req.Body).Decode(&m)
		var user_id string = m.User_id

		params := mux.Vars(req)
		var photo_id string = params["photo_id"]
		fmt.Println(photo_id)
		fmt.Println(user_id)
		createLike(photo_id, user_id)
		like_count := readLikes(photo_id)
		// Call the SNS
		go publishSNS(photo_id, like_count, LIKES_SNS_TOPIC)
		formatter.JSON(w, http.StatusOK, struct{ Status string }{"Liked"})
	}
}


// API to DLETE a comment
func removeLike(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var m like
		_ = json.NewDecoder(req.Body).Decode(&m)
		var user_id string = m.User_id

		params := mux.Vars(req)
		var photo_id string = params["photo_id"]
		fmt.Println(photo_id)
		fmt.Println(user_id)
		// Call the SNS
		// Update mongodb (delete like for photo and backup)
		formatter.JSON(w, http.StatusOK, struct{ Status string }{"Like Removed"})
	}
}

// API to get comment list for a photo
func commentList(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		enableCors(&w)
		params := mux.Vars(req)
		var photo_id string = params["photo_id"]
		var user_id string = params["user_id"]
		fmt.Println(photo_id)
		fmt.Println(user_id)
		comments := readComments(photo_id, user_id)
		// Code to get the value from database
		formatter.JSON(w, http.StatusOK, comments)
	}
}


// API to POST a new comment
func addComment(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		setupResponse(&w, req)
		if (*req).Method == "OPTIONS" {
			return
		}
		var m user_comment
		_ = json.NewDecoder(req.Body).Decode(&m)
		var user_id string = m.User_id
		var comment string = m.Comment
		var user_name string = m.User_name

		params := mux.Vars(req)
		var photo_id string = params["photo_id"]
		fmt.Println(photo_id)
		fmt.Println(user_name)
		fmt.Println(user_id)
		fmt.Println(comment)

		createComment(photo_id, user_id, user_name, comment)
		comment_count := readCommentCount(photo_id)
		//call SNS
		go publishSNS(photo_id, comment_count, COMMENTS_SNS_TOPIC)
		formatter.JSON(w, http.StatusOK, struct{ Status string }{"Comment Added"})
	}
}

// API to DLETE a comment
func removeComment(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var m user_comment
		_ = json.NewDecoder(req.Body).Decode(&m)
		var user_id string = m.User_id
		var comment string = m.Comment

		params := mux.Vars(req)
		var photo_id string = params["photo_id"]
		fmt.Println(photo_id)
		fmt.Println(user_id)
		fmt.Println(comment)
		// Call the SNS
		// Update mongodb (delete like for photo and backup)
		deleteComment(photo_id, user_id, comment)
		formatter.JSON(w, http.StatusOK, struct{ Status string }{"Comment Removed"})
	}
}

// API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"Comment and Like API Version 1.0"})
	}
}