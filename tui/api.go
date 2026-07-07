package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

type Movie struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Director    string `json:"director"`
	Genre       string `json:"genre"`
	ReleaseYear int    `json:"release_year"`
	Description string `json:"description"`
}

func loginCmd(email, password string) tea.Cmd {
	return func() tea.Msg {
		body := map[string]string{
			"email":    email,
			"password": password,
		}
		b, _ := json.Marshal(body)

		resp, err := http.Post(API_URL+"/user/login", "application/json", bytes.NewBuffer(b))
		if err != nil {
			return loginMsg{err: err}
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return loginMsg{err: fmt.Errorf("status %d: %s", resp.StatusCode, string(data))}
		}

		var res map[string]string
		if err := json.Unmarshal(data, &res); err != nil {
			return loginMsg{err: err}
		}

		token, ok := res["token"]
		if !ok {
			return loginMsg{err: fmt.Errorf("token manquant")}
		}

		return loginMsg{token: token}
	}
}

func getMoviesCmd(token string) tea.Cmd {
	return func() tea.Msg {
		req, _ := http.NewRequest("GET", API_URL+"/movies", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return moviesMsg{err: err}
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return moviesMsg{err: fmt.Errorf("status %d: %s", resp.StatusCode, string(data))}
		}

		var movies []Movie
		if err := json.Unmarshal(data, &movies); err != nil {
			return moviesMsg{err: err}
		}

		return moviesMsg{movies: movies}
	}
}

func searchMoviesCmd(token, query string) tea.Cmd {
	return func() tea.Msg {
		req, _ := http.NewRequest("GET", API_URL+"/search?query="+query, nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return moviesMsg{err: err}
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return moviesMsg{err: fmt.Errorf("status %d: %s", resp.StatusCode, string(data))}
		}

		var movies []Movie
		if err := json.Unmarshal(data, &movies); err != nil {
			return moviesMsg{err: err}
		}

		return moviesMsg{movies: movies}
	}
}

func createMovieCmd(token, title, director, genre, year, desc string) tea.Cmd {
	return func() tea.Msg {

		yearInt, _ := strconv.Atoi(year)

		body := map[string]interface{}{
			"title":        title,
			"director":     director,
			"genre":        genre,
			"release_year": yearInt,
			"description":  desc,
		}

		b, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", API_URL+"/movies", bytes.NewBuffer(b))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return createMsg{err: err}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			data, _ := io.ReadAll(resp.Body)
			return createMsg{err: fmt.Errorf("status %d: %s", resp.StatusCode, string(data))}
		}

		return createMsg{err: nil}
	}
}
