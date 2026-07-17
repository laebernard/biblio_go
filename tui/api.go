package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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

func registerCmd(name, email, password string) tea.Cmd {
	return func() tea.Msg {
		body := map[string]string{
			"name":     name,
			"email":    email,
			"password": password,
		}
		b, _ := json.Marshal(body)

		resp, err := http.Post(API_URL+"/user/register", "application/json", bytes.NewBuffer(b))
		if err != nil {
			return registerMsg{err: err}
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return registerMsg{err: fmt.Errorf("status %d: %s", resp.StatusCode, string(data))}
		}

		var res map[string]string
		if err := json.Unmarshal(data, &res); err != nil {
			return registerMsg{err: err}
		}

		token, ok := res["token"]
		if !ok {
			return registerMsg{err: fmt.Errorf("token manquant")}
		}

		return registerMsg{token: token}
	}
}

func updateOwnProfileCmd(token, name, email, password string) tea.Cmd {
	return func() tea.Msg {
		body := map[string]interface{}{}
		if name != "" {
			body["name"] = name
		}
		if email != "" {
			body["email"] = email
		}
		if password != "" {
			body["password"] = password
		}

		b, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", API_URL+"/user/me", bytes.NewBuffer(b))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return userUpdateMsg{err: err}
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return userUpdateMsg{err: fmt.Errorf("status %d: %s", resp.StatusCode, string(data))}
		}

		return userUpdateMsg{err: nil}
	}
}

func updateUserByAdminCmd(token, id, name, email, password string) tea.Cmd {
	return func() tea.Msg {
		body := map[string]interface{}{}
		if name != "" {
			body["name"] = name
		}
		if email != "" {
			body["email"] = email
		}
		if password != "" {
			body["password"] = password
		}

		b, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", API_URL+"/user/"+id, bytes.NewBuffer(b))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return adminUserUpdateMsg{err: err}
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return adminUserUpdateMsg{err: fmt.Errorf("status %d: %s", resp.StatusCode, string(data))}
		}

		return adminUserUpdateMsg{err: nil}
	}
}

func resetDatabaseCmd(token string) tea.Cmd {
	return func() tea.Msg {
		req, _ := http.NewRequest("DELETE", API_URL+"/reset", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return resetMsg{err: err}
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return resetMsg{err: fmt.Errorf("status %d: %s", resp.StatusCode, string(data))}
		}

		return resetMsg{err: nil}
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

func getMovieCmd(token, id string) tea.Cmd {
	return func() tea.Msg {
		req, _ := http.NewRequest("GET", API_URL+"/movies/"+id, nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return movieMsg{err: err}
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return movieMsg{err: fmt.Errorf("status %d: %s", resp.StatusCode, string(data))}
		}

		var movie Movie
		if err := json.Unmarshal(data, &movie); err != nil {
			return movieMsg{err: err}
		}

		return movieMsg{movie: movie}
	}
}

func searchMoviesCmd(token, query string) tea.Cmd {
	return func() tea.Msg {
		req, _ := http.NewRequest("GET", API_URL+"/search?query="+url.QueryEscape(query), nil)
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

func updateMovieCmd(token, id, title, director, genre, year, desc string) tea.Cmd {
	return func() tea.Msg {
		if strings.TrimSpace(title) == "" &&
			strings.TrimSpace(director) == "" &&
			strings.TrimSpace(genre) == "" &&
			strings.TrimSpace(year) == "" &&
			strings.TrimSpace(desc) == "" {
			return updateMsg{err: fmt.Errorf("aucun champ à mettre à jour")}
		}

		getReq, _ := http.NewRequest("GET", API_URL+"/movies/"+id, nil)
		getReq.Header.Set("Authorization", "Bearer "+token)

		getResp, err := http.DefaultClient.Do(getReq)
		if err != nil {
			return updateMsg{err: err}
		}
		defer getResp.Body.Close()

		if getResp.StatusCode != http.StatusOK {
			data, _ := io.ReadAll(getResp.Body)
			return updateMsg{err: fmt.Errorf("status %d: %s", getResp.StatusCode, string(data))}
		}

		var current Movie
		if err := json.NewDecoder(getResp.Body).Decode(&current); err != nil {
			return updateMsg{err: err}
		}

		if strings.TrimSpace(title) != "" {
			current.Title = title
		}
		if strings.TrimSpace(director) != "" {
			current.Director = director
		}
		if strings.TrimSpace(genre) != "" {
			current.Genre = genre
		}
		if strings.TrimSpace(desc) != "" {
			current.Description = desc
		}
		if strings.TrimSpace(year) != "" {
			yearInt, err := strconv.Atoi(strings.TrimSpace(year))
			if err != nil || yearInt <= 0 {
				return updateMsg{err: fmt.Errorf("année invalide")}
			}
			current.ReleaseYear = yearInt
		}

		body := map[string]interface{}{
			"title":        current.Title,
			"director":     current.Director,
			"genre":        current.Genre,
			"release_year": current.ReleaseYear,
			"description":  current.Description,
		}

		b, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", API_URL+"/movies/"+id, bytes.NewBuffer(b))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return updateMsg{err: err}
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return updateMsg{err: fmt.Errorf("status %d: %s", resp.StatusCode, string(data))}
		}

		return updateMsg{err: nil}
	}
}

func deleteMovieCmd(token, id string) tea.Cmd {
	return func() tea.Msg {
		req, _ := http.NewRequest("DELETE", API_URL+"/movies/"+id, nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return deleteMsg{err: err}
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return deleteMsg{err: fmt.Errorf("status %d: %s", resp.StatusCode, string(data))}
		}

		return deleteMsg{err: nil}
	}
}

func bulkDeleteMoviesCmd(token string, ids []uint) tea.Cmd {
	return func() tea.Msg {
		payload := make([]map[string]uint, len(ids))
		for i, id := range ids {
			payload[i] = map[string]uint{"id": id}
		}

		b, _ := json.Marshal(payload)
		req, _ := http.NewRequest("DELETE", API_URL+"/movies", bytes.NewBuffer(b))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return bulkDeleteMsg{err: err}
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return bulkDeleteMsg{err: fmt.Errorf("status %d: %s", resp.StatusCode, string(data))}
		}

		return bulkDeleteMsg{err: nil}
	}
}

func parseIDList(input string) ([]uint, error) {
	parts := strings.Split(input, ",")
	ids := make([]uint, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}

		id, err := strconv.ParseUint(trimmed, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("id invalide %q: %w", trimmed, err)
		}
		ids = append(ids, uint(id))
	}

	if len(ids) == 0 {
		return nil, fmt.Errorf("aucun id fourni")
	}

	return ids, nil
}

func buildBulkUpdateMovies(ids []uint, title, director, genre, year, desc string) ([]Movie, error) {
	yearInt := 0
	if year != "" {
		var err error
		yearInt, err = strconv.Atoi(year)
		if err != nil {
			return nil, fmt.Errorf("année invalide: %w", err)
		}
	}

	movies := make([]Movie, 0, len(ids))
	for _, id := range ids {
		movies = append(movies, Movie{
			ID:          id,
			Title:       title,
			Director:    director,
			Genre:       genre,
			ReleaseYear: yearInt,
			Description: desc,
		})
	}

	return movies, nil
}

func bulkUpdateMoviesCmd(token string, ids []uint, title, director, genre, year, desc string) tea.Cmd {
	return func() tea.Msg {
		movies, err := buildBulkUpdateMovies(ids, title, director, genre, year, desc)
		if err != nil {
			return bulkUpdateMsg{err: err}
		}

		b, _ := json.Marshal(movies)
		req, _ := http.NewRequest("PUT", API_URL+"/movies/bulk", bytes.NewBuffer(b))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return bulkUpdateMsg{err: err}
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return bulkUpdateMsg{err: fmt.Errorf("status %d: %s", resp.StatusCode, string(data))}
		}

		return bulkUpdateMsg{err: nil}
	}
}
