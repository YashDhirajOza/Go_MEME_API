package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Meme struct {
	PostLink  string   `json:"postLink"`
	Subreddit string   `json:"subreddit"`
	Title     string   `json:"title"`
	URL       string   `json:"url"`
	NSFW      bool     `json:"nsfw"`
	Spoiler   bool     `json:"spoiler"`
	Author    string   `json:"author"`
	Ups       int      `json:"ups"`
	Preview   []string `json:"preview"`
}

func FetchMeme() (*Meme, error) {
	allowedSubreddits := []string{"dankmemes", "memes", "wholesomememes"}
	apiURL := fmt.Sprintf("https://meme-api.com/gimme/%s", strings.Join(allowedSubreddits, ","))

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var meme Meme
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("API Response:", string(body))
	if err := json.Unmarshal(body, &meme); err != nil {
		return nil, err
	}

	return &meme, nil
}

func MemeHandler(w http.ResponseWriter, r *http.Request) {
	meme, err := FetchMeme()
	if err != nil {
		http.Error(w, "Error fetching meme: "+err.Error(), http.StatusInternalServerError)
		return
	}

	html := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Random Meme</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				text-align: center;
				background-color: #282c34;
				color: white;
				margin: 0;
				padding: 0;
			}
			img {
				max-width: 500px;
				width: 90%;
				height: auto;
				border-radius: 15px;
				margin-top: 20px;
				box-shadow: 0 4px 10px rgba(0, 0, 0, 0.3);
			}
			a {
				color: #61dafb;
				text-decoration: none;
			}
		</style>
	</head>
	<body>
		<h1>` + meme.Title + `</h1>
		<p>Subreddit: <strong>` + meme.Subreddit + `</strong></p>
		<p>Author: <strong>` + meme.Author + `</strong></p>
		<p>Upvotes: <strong>` + fmt.Sprintf("%d", meme.Ups) + `</strong></p>
		<a href="` + meme.PostLink + `" target="_blank">View Post on Reddit</a>
		<br><br>
		<img src="` + meme.URL + `" alt="Meme Image">
	</body>
	</html>
	`

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func main() {
	http.HandleFunc("/random-meme", MemeHandler)

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
