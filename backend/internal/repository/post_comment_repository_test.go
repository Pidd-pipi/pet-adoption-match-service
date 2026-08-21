package repository

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/gbadopt/gbadopt/internal/model"
)

func TestPostCommentRepoPreservesError(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("gorm open: %v", err)
	}
	repo := NewPostCommentRepository(db)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "post_comments"`)).WillReturnError(errors.New("db write failed"))
	if err := repo.Create(&model.PostComment{PostID: 1, UserID: 1, Content: "hi"}); err == nil {
		t.Fatalf("expected create error")
	}

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "post_comments"`)).WillReturnError(errors.New("db delete failed"))
	if err := repo.Delete(1); err == nil {
		t.Fatalf("expected delete error")
	}
}
