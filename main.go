package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
	"os"

	"github.com/Warren-Wang-OG/go-social-media-backend/database"
)

type errorBody struct {
	Error string `json:"error"`
}

type apiConfig struct {
	dbClient database.Client
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK,
		database.User{
			Email: "test@example.com",
		})
}

func testErrHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, errors.New("error handler default response"))
}

// wrapper for respondWithJSON for sending errors as the interface used to be converted to json
func respondWithError(w http.ResponseWriter, code int, err error) {
	respondWithJSON(w, code, errorBody{Error: err.Error()})
}

// handles http requests and return json
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(response)
	w.WriteHeader(code)
}

// TODO: this function is not used by any other func, except the test function as of right now
// a filter for eligible users
func userIsEligible(email, password string, age int) error {
	// empty email or password
	if email == "" {
		return errors.New("email can't be empty")
	}
	if password == "" {
		return errors.New("password can't be empty")
	}

	// age is less than 18
	if age < 18 {
		return fmt.Errorf("age %d is less than 18, must be at least 18", age)
	}

	return nil
}

// TODO: allow updating email? would need to delete old user, then create new
// key-value map, but then Posts would be affected to, since those have an email attached to them
// could get all posts from the user via the old email, then update all of those posts to use the new email
// then write those changes to disk
//
// update user when a PUT request is made to /users/EMAIL
// take parameters in the body of the request
func (apiCfg apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Name     string `json:"name"`
		Age      int    `json:"age"`
	}

	// get email from path
	email := r.URL.Path[len("/users/"):]
	// get parameters from body
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	// update user
	_, err = apiCfg.dbClient.UpdateUser(email, params.Password, params.Name, params.Age)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	// good, return 200 status code
	respondWithJSON(w, http.StatusOK, struct{}{})
}

// return a user when a GET request is made to /users/EMAIL
// return marshalled json of the user
func (apiCfg apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	// get email from path
	email := r.URL.Path[len("/users/"):]
	user, err := apiCfg.dbClient.GetUser(email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	// good, return 200 status code with the user info
	respondWithJSON(w, http.StatusOK, user)
}

// delete a user when a DELETE request is made to /users/EMAIL
// conventionally, DELETE requests take no body
// so we find user based on email in the path /users/EMAIL
func (apiCfg apiConfig) handlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	// get email from path
	email := r.URL.Path[len("/users/"):]

	// delete user
	err := apiCfg.dbClient.DeleteUser(email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	// good, return 200 status code
	respondWithJSON(w, http.StatusOK, struct{}{})
}

// create user when a POST request is made to /users
// this takes an input of a json object with the following fields:
// email, password, name, age
func (apiCfg apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Age      int    `json:"age"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	// create the new user from params
	_, err = apiCfg.dbClient.CreateUser(params.Email, params.Password, params.Name, params.Age)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	// good, return 201 status code
	respondWithJSON(w, http.StatusCreated, struct{}{})
}

// routes the specific http request about users to the correct handler
func (apiCfg apiConfig) endpointUsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// call GET handler
		apiCfg.handlerGetUser(w, r)
	case http.MethodPost:
		// call POST handler
		apiCfg.handlerCreateUser(w, r)
	case http.MethodPut:
		// call PUT handler
		apiCfg.handlerUpdateUser(w, r)
	case http.MethodDelete:
		// call DELETE handler
		apiCfg.handlerDeleteUser(w, r)
	default:
		respondWithError(w, 404, errors.New("method not supported"))
	}
}

// TODO: create a post update function, may require changes to the 
// original post definition.

// delete a post when a DELETE request is made to /posts/ID
// ID is UUID of the post to be deleted
func (apiCfg apiConfig) handlerDeletePost(w http.ResponseWriter, r *http.Request) {
	// get uuid from path
	uuid := r.URL.Path[len("/posts/"):]

	// delete post
	err := apiCfg.dbClient.DeletePost(uuid)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	// good, return 200 status code
	respondWithJSON(w, http.StatusOK, struct{}{})
}

// create a post when a POST request is made to /posts
// this takes an input of a json object with the following fields:
// email, text
func (apiCfg apiConfig) handlerCreatePost(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserEmail string `json:"email"`
		Text      string `json:"text"`
	}

	// convert json object to parameters struct
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	// create the new post from params
	_, err = apiCfg.dbClient.CreatePost(params.UserEmail, params.Text)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	// good, return 201 status code
	respondWithJSON(w, http.StatusCreated, struct{}{})
}

// return list of all posts when a GET request is made to /posts/EMAIL
// for a specific user based on the EMAIL given
func (apiCfg apiConfig) handlerRetrievePosts(w http.ResponseWriter, r *http.Request) {
	// get email from path
	email := r.URL.Path[len("/posts/"):]
	posts, err := apiCfg.dbClient.GetPosts(email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	// good, return 200 status code with the user info
	respondWithJSON(w, http.StatusOK, posts)
}

// routes specific http request about posts to the correct handler
func (apiCfg apiConfig) endpointPostsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// call GET handler
		apiCfg.handlerRetrievePosts(w, r)
	case http.MethodPost:
		// call POST handler
		apiCfg.handlerCreatePost(w, r)
	case http.MethodDelete:
		// call DELETE handler
		apiCfg.handlerDeletePost(w, r)
	default:
		respondWithError(w, 404, errors.New("method not supported"))
	}
}

func main() {
	// create a new database
	c := database.NewClient("db.json")
	err := c.EnsureDB()
	if err != nil {
		log.Fatal(err)
	}

	// create api construct to work with handlers instead of directly
	// with the client
	apiConfig := apiConfig{
		dbClient: c,
	}

	// allocate http request multiplexer
	serveMux := http.NewServeMux()

	// handler to register at the "/" root path
	serveMux.HandleFunc("/", testHandler)

	// handler to register at the "/error" path
	serveMux.HandleFunc("/err", testErrHandler)

	// handler to register at the "/users" path
	serveMux.HandleFunc("/users", apiConfig.endpointUsersHandler)
	serveMux.HandleFunc("/users/", apiConfig.endpointUsersHandler)

	// Posts
	serveMux.HandleFunc("/posts", apiConfig.endpointPostsHandler)
	serveMux.HandleFunc("/posts/", apiConfig.endpointPostsHandler)

	// http server
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%v", port)
	fmt.Printf("Running on port: %v\n", port)
	srv := http.Server{
		Handler:      serveMux,
		Addr:         addr,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	// wait and listen
	srv.ListenAndServe()
}
