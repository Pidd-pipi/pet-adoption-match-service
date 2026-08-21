package service

import (
	"testing"

	"github.com/gbadopt/gbadopt/internal/model"
)

func TestStatsNilMapSafety(t *testing.T) {
	counts := OrgAdoptionCounts([]model.AdoptionApplication{
		{OrgID: 1, Status: "approved"},
		{OrgID: 1, Status: "submitted"},
		{OrgID: 2, Status: "approved"},
	})
	if counts[1]["approved"] != 1 || counts[1]["submitted"] != 1 || counts[2]["approved"] != 1 {
		t.Fatalf("unexpected counts: %+v", counts)
	}
	species := SpeciesCounts([]model.Pet{
		{Species: "dog"},
		{Species: "dog"},
		{Species: "cat"},
	})
	if species["dog"] != 2 || species["cat"] != 1 {
		t.Fatalf("unexpected species: %+v", species)
	}
}
