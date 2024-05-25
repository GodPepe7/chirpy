package db

import (
	"testing"
)

func TestCreateChirp(t *testing.T) {
	RemoveDBFile()
	defer RemoveDBFile()
	database, err := NewDB("../../database")
	if err != nil {
		t.Errorf("expected no errors: %v", err)
		return
	}
	chirp, err := database.CreateChirp("test", 1)
	if err != nil {
		t.Errorf("expected no errors: %v", err)
		return
	}
	if chirp.Body != "test" {
		t.Errorf("expected body to be 'test' instead of: %v", chirp.Body)
		return
	}
	if chirp.Id != 1 {
		t.Errorf("expected id to be '1' instead of: %v", chirp.Id)
		return
	}

	chirps, err := database.GetChirps()
	if err != nil {
		t.Errorf("expected no errors: %v", err)
		return
	}
	if len(chirps) != 1 {
		t.Errorf("expected exactly one chirp instead of: %v", len(chirps))
		return
	}
	if chirps[0].Id != 1 && chirps[0].Body != "test" {
		t.Errorf("expected received chirp to be identical to created one: %v, %v", chirps[0], chirp)
		return
	}
}
