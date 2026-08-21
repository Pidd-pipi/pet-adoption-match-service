package service

import (
	"errors"
	"regexp"
	"testing"

	"gorm.io/gorm"

	"github.com/gbadopt/gbadopt/internal/repository"
	"github.com/gbadopt/gbadopt/internal/util"
)

func TestDonationErrorChainAndMapping(t *testing.T) {
	t.Run("usage missing donation maps to 404", func(t *testing.T) {
		db, mock := newServiceDB(t)
		svc := NewDonationService(repository.NewDonationRepository(db), repository.NewDonationUsageRepository(db), repository.NewOrganizationRepository(db), newTestLogger())
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "donations" WHERE "donations"."id" = $1 ORDER BY "donations"."id" LIMIT $2`)).
			WithArgs(999, 1).WillReturnError(gorm.ErrRecordNotFound)
		_, err := svc.CreateUsage(1, 999, 50, "desc", "")
		var appErr *util.AppError
		if !errors.As(err, &appErr) || appErr.HTTPStatus != 404 {
			t.Fatalf("expected 404 AppError, got %v", err)
		}
	})
	t.Run("donate system error not mapped to 404", func(t *testing.T) {
		db, mock := newServiceDB(t)
		svc := NewDonationService(repository.NewDonationRepository(db), repository.NewDonationUsageRepository(db), repository.NewOrganizationRepository(db), newTestLogger())
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE "organizations"."id" = $1 ORDER BY "organizations"."id" LIMIT $2`)).
			WithArgs(999, 1).WillReturnError(errors.New("db down"))
		_, err := svc.Donate(1, 999, 100)
		var appErr *util.AppError
		if errors.As(err, &appErr) && appErr.HTTPStatus == 404 {
			t.Fatalf("system error should not map to 404, got %v", err)
		}
	})
}
