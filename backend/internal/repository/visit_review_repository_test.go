package repository

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestMarkOverdueOnlyPending(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("gorm open: %v", err)
	}
	repo := NewVisitReviewRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "visit_reviews" SET "status"=$1 WHERE user_id = $2 AND status = $3 AND due_date < $4`)).
		WithArgs("overdue", 1, "pending", sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	if err := repo.MarkOverdue(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListByUserAllStatuses(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("gorm open: %v", err)
	}
	repo := NewVisitReviewRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "visit_reviews" WHERE user_id = $1 ORDER BY due_date ASC`)).
		WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "application_id", "user_id", "org_id", "scheduled_days", "due_date", "status", "photos", "note", "created_at"}))

	if _, err := repo.ListByUser(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}


func TestListByOrgAllStatuses(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("gorm open: %v", err)
	}
	repo := NewVisitReviewRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "visit_reviews" WHERE org_id = $1 ORDER BY due_date ASC`)).
		WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "application_id", "user_id", "org_id", "scheduled_days", "due_date", "status", "photos", "note", "created_at"}))

	if _, err := repo.ListByOrg(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
