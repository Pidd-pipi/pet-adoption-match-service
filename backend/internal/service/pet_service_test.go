package service

import (
	"context"
	"errors"
	"testing"

	"github.com/gbadopt/gbadopt/internal/repository"
)

func TestPetListContextCancellation(t *testing.T) {
	db, _ := newServiceDB(t)
	repo := repository.NewPetRepository(db)
	svc := NewPetService(repo, nil, nil, newTestLogger())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, _, err := svc.List(ctx, "dog", "", "", "", 1, 10)
	if err == nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
