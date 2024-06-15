package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Movie struct {
	ID       string    `json:"id"`
	Isdn     string    `json:"isdn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

// route to get all movies
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// route to delete a movie by id
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Get the movie ID from the URL parameters
	id := chi.URLParam(r, "id")
	for index, item := range movies {
		if item.ID == id {
			movies = append(movies[:index], movies[index+1:]...)
			json.NewEncoder(w).Encode(movies)
			return
		}
	}
}

// route to get a movie by id
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Get the movie ID from the URL parameters
	id := chi.URLParam(r, "id")
	for _, item := range movies {
		if item.ID == id {

			json.NewEncoder(w).Encode(item)
			return
		}
	}
}
func createMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(100))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movies)

}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode the JSON request body into a Movie struct
	var updatedMovie Movie
	if err := json.NewDecoder(r.Body).Decode(&updatedMovie); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the movie ID from the URL parameters
	id := chi.URLParam(r, "id")

	// Find the index of the movie with the specified ID in the movies slice
	index := -1
	for i, item := range movies {
		if item.ID == id {
			index = i
			break
		}
	}

	// If the movie with the specified ID is not found, return an error
	if index == -1 {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}

	// Update the movie with the new details
	movies[index] = updatedMovie
	movies[index].ID = id

	// Return the updated movie as JSON
	json.NewEncoder(w).Encode(movies)
}
func main() {
	r := chi.NewRouter()
	movies = []Movie{
		{ID: "1", Isdn: "1234567890", Title: "Inception", Director: &Director{Firstname: "Christopher", Lastname: "Nolan"}},
		{ID: "2", Isdn: "9876543210", Title: "The Shawshank Redemption", Director: &Director{Firstname: "Frank", Lastname: "Darabont"}},
	}
	r.Get("/movies", getMovies)
	r.Get("/movie/{id}", getMovie)
	r.Post("/createmovie", createMovies)
	r.Put("/updatemovie/{id}", updateMovie)
	r.Delete("/deletemovie/{id}", deleteMovie)

	fmt.Println("starting server at port:3050")

	log.Fatal(http.ListenAndServe(":3050", r))
}
