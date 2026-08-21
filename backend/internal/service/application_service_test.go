package service

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/gbadopt/gbadopt/internal/repository"
	"github.com/gbadopt/gbadopt/internal/util"
)

func TestApplicationStatusMachine(t *testing.T) {
	db, mock := newServiceDB(t)
	svc := NewApplicationService(db, repository.NewAdoptionApplicationRepository(db), repository.NewPetRepository(db), repository.NewOrganizationRepository(db), newTestLogger())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "adoption_applications" WHERE "adoption_applications"."id" = $1 ORDER BY "adoption_applications"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "pet_id", "org_id", "questionnaire", "status", "created_at", "updated_at"}).
			AddRow(1, 2, 1, 1, "{}", "offline_interview", time.Now(), time.Now()))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE user_id = $1 ORDER BY "organizations"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "cert_type", "status", "contact", "city", "description", "created_at"}).
			AddRow(1, 1, "org", "registered", "approved", "138", "上海", "desc", time.Now()))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "pets" WHERE "pets"."id" = $1 ORDER BY "pets"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "org_id", "name", "species", "breed", "age", "gender", "size", "city", "description", "personality", "health_status", "neutered", "vaccinated", "image_urls", "status", "created_at"}).
			AddRow(1, 1, "旺财", "dog", "田园犬", 2, "male", "medium", "上海", "d", "p", "h", true, true, "[]", "pending", time.Now()))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "adoption_applications"`)).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "pets" SET`)).WithArgs(
		sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "adopted",
		sqlmock.AnyArg(), 1,
	).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	_, err := svc.UpdateStatus(1, 1, "org", "approved")
	var appErr *util.AppError
	if errors.As(err, &appErr) && appErr.HTTPStatus == 403 {
		t.Fatalf("org owner should not get 403")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}


func TestApplicationListByOrgStatus(t *testing.T) {
	db, mock := newServiceDB(t)
	svc := NewApplicationService(db, repository.NewAdoptionApplicationRepository(db), repository.NewPetRepository(db), repository.NewOrganizationRepository(db), newTestLogger())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE user_id = $1 ORDER BY "organizations"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "cert_type", "status", "contact", "city", "description", "created_at"}).
			AddRow(1, 1, "org", "registered", "approved", "138", "上海", "desc", time.Now()))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "adoption_applications" WHERE org_id = $1 AND status = $2 ORDER BY id DESC`)).
		WithArgs(1, "submitted").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "pet_id", "org_id", "questionnaire", "status", "created_at", "updated_at"}))

	if _, err := svc.ListByOrg(1, "submitted"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApplicationRejectRevertsPet(t *testing.T) {
	db, mock := newServiceDB(t)
	svc := NewApplicationService(db, repository.NewAdoptionApplicationRepository(db), repository.NewPetRepository(db), repository.NewOrganizationRepository(db), newTestLogger())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "adoption_applications" WHERE "adoption_applications"."id" = $1 ORDER BY "adoption_applications"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "pet_id", "org_id", "questionnaire", "status", "created_at", "updated_at"}).
			AddRow(1, 2, 1, 1, "{}", "communicating", time.Now(), time.Now()))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE user_id = $1 ORDER BY "organizations"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "cert_type", "status", "contact", "city", "description", "created_at"}).
			AddRow(1, 1, "org", "registered", "approved", "138", "上海", "desc", time.Now()))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "adoption_applications"`)).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "pets" WHERE "pets"."id" = $1 ORDER BY "pets"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "org_id", "name", "species", "breed", "age", "gender", "size", "city", "description", "personality", "health_status", "neutered", "vaccinated", "image_urls", "status", "created_at"}).
			AddRow(1, 1, "旺财", "dog", "田园犬", 2, "male", "medium", "上海", "d", "p", "h", true, true, "[]", "pending", time.Now()))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "pets" SET`)).WithArgs(
		sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "available",
		sqlmock.AnyArg(), 1,
	).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	if _, err := svc.UpdateStatus(1, 1, "org", "rejected"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("pet was not reverted to available: %v", err)
	}
}
