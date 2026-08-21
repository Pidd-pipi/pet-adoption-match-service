package service

import (
	"errors"
	"regexp"
	"testing"

	"gorm.io/gorm"

	"github.com/gbadopt/gbadopt/internal/repository"
	"github.com/gbadopt/gbadopt/internal/util"
)

func TestOrgGetMissingNoPanic(t *testing.T) {
	db, mock := newServiceDB(t)
	svc := NewOrganizationService(repository.NewOrganizationRepository(db), newTestLogger())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE "organizations"."id" = $1 ORDER BY "organizations"."id" LIMIT $2`)).
		WithArgs(999, 1).WillReturnError(gorm.ErrRecordNotFound)

	_, err := svc.Get(999)
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.HTTPStatus != 404 {
		t.Fatalf("expected 404, got %v", err)
	}
}

func TestOrgReviewMissingNoPanic(t *testing.T) {
	db, mock := newServiceDB(t)
	svc := NewOrganizationService(repository.NewOrganizationRepository(db), newTestLogger())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE "organizations"."id" = $1 ORDER BY "organizations"."id" LIMIT $2`)).
		WithArgs(999, 1).WillReturnError(gorm.ErrRecordNotFound)

	_, err := svc.Review(999, "approved")
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.HTTPStatus != 404 {
		t.Fatalf("expected 404, got %v", err)
	}
}
