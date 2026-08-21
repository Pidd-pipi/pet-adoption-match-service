package handler

import (
	"testing"

	"github.com/gbadopt/gbadopt/internal/model"
)

func TestHomeHelpersNilMapSafety(t *testing.T) {
	species := summarizeSpecies([]model.Pet{
		{Species: "dog"},
		{Species: "rabbit"},
	})
	if species["dog"] != 1 || species["rabbit"] != 1 {
		t.Fatalf("unexpected species: %+v", species)
	}
	cities := summarizeCities([]model.Organization{
		{City: "上海"},
		{City: "上海"},
		{City: "北京"},
	})
	if cities["上海"] != 2 || cities["北京"] != 1 {
		t.Fatalf("unexpected cities: %+v", cities)
	}
}
