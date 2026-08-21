package service

import (
	"errors"
	"regexp"
	"testing"

	"gorm.io/gorm"

	"github.com/gbadopt/gbadopt/internal/repository"
	"github.com/gbadopt/gbadopt/internal/util"
)

func TestFavoriteAddDuplicateConflict(t *testing.T) {
	db, mock := newServiceDB(t)
	svc := NewFavoriteService(repository.NewFavoriteRepository(db), newTestLogger())

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "favorites"`)).WillReturnError(errors.New("duplicate key value violates unique constraint"))

	_, err := svc.Add(1, "pet", 1)
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.HTTPStatus != 409 {
		t.Fatalf("expected 409 conflict, got %v", err)
	}
}

func TestFavoriteRemoveMissing(t *testing.T) {
	db, mock := newServiceDB(t)
	svc := NewFavoriteService(repository.NewFavoriteRepository(db), newTestLogger())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "favorites" WHERE user_id = $1 AND target_type = $2 AND target_id = $3 ORDER BY "favorites"."id" LIMIT $4`)).
		WithArgs(1, "pet", 1, 1).WillReturnError(gorm.ErrRecordNotFound)

	err := svc.Remove(1, "pet", 1)
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.HTTPStatus != 404 {
		t.Fatalf("expected 404 not found, got %v", err)
	}
}
