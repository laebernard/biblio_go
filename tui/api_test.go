package main

import (
	"reflect"
	"testing"
)

func TestParseIDList(t *testing.T) {
	ids, err := parseIDList("1, 2, 3")
	if err != nil {
		t.Fatalf("parseIDList returned error: %v", err)
	}

	expected := []uint{1, 2, 3}
	if !reflect.DeepEqual(ids, expected) {
		t.Fatalf("expected %v, got %v", expected, ids)
	}
}

func TestBuildBulkUpdateMovies(t *testing.T) {
	movies, err := buildBulkUpdateMovies([]uint{1, 2}, "Titre", "Dir", "Genre", "1999", "Desc")
	if err != nil {
		t.Fatalf("buildBulkUpdateMovies returned error: %v", err)
	}

	if len(movies) != 2 {
		t.Fatalf("expected 2 movies, got %d", len(movies))
	}

	if movies[0].ID != 1 || movies[1].ID != 2 {
		t.Fatalf("unexpected ids: %+v", movies)
	}

	if movies[0].Title != "Titre" || movies[0].Director != "Dir" || movies[0].Genre != "Genre" || movies[0].ReleaseYear != 1999 || movies[0].Description != "Desc" {
		t.Fatalf("unexpected movie payload: %+v", movies[0])
	}
}
