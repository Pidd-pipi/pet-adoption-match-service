package service

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/gbadopt/gbadopt/internal/model"
	"github.com/gbadopt/gbadopt/internal/repository"
)

func TestReviewOverdueSubmit(t *testing.T) {
	db, mock := newServiceDB(t)
	svc := NewReviewService(repository.NewVisitReviewRepository(db), repository.NewAdoptionApplicationRepository(db), repository.NewOrganizationRepository(db), newTestLogger())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "visit_reviews" WHERE "visit_reviews"."id" = $1 ORDER BY "visit_reviews"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "application_id", "user_id", "org_id", "scheduled_days", "due_date", "status", "photos", "note", "created_at"}).
			AddRow(1, 1, 1, 1, 30, time.Now(), model.ReviewOverdue, "[]", "", time.Now()))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "visit_reviews"`)).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	v, err := svc.Submit(1, 1, "[]", "note")
	if err != nil {
		t.Fatalf("overdue review should be submittable: %v", err)
	}
	if v.Status != model.ReviewSubmitted {
		t.Fatalf("expected submitted status, got %s", v.Status)
	}
}


func TestReviewCreateDefaultDueDate(t *testing.T) {
	db, mock := newServiceDB(t)
	svc := NewReviewService(repository.NewVisitReviewRepository(db), repository.NewAdoptionApplicationRepository(db), repository.NewOrganizationRepository(db), newTestLogger())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "adoption_applications" WHERE "adoption_applications"."id" = $1 ORDER BY "adoption_applications"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "pet_id", "org_id", "questionnaire", "status", "created_at", "updated_at"}).
			AddRow(1, 1, 1, 1, "{}", "approved", time.Now(), time.Now()))
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "visit_reviews"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	v, err := svc.Create(1, 1, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.ScheduledDays != 30 {
		t.Fatalf("expected default 30 days, got %d", v.ScheduledDays)
	}
	if !v.DueDate.After(time.Now()) {
		t.Fatalf("due date should be in the future, got %v", v.DueDate)
	}
}
