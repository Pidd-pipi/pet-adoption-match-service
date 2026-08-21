package service

import (
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/gbadopt/gbadopt/internal/repository"
)

func newRows(cols ...string) *sqlmock.Rows {
	return sqlmock.NewRows(cols)
}

func TestCommentCreatePreservesError(t *testing.T) {
	db, mock := newServiceDB(t)
	svc := NewCommentService(db, repository.NewPostCommentRepository(db), repository.NewCommunityPostRepository(db), newTestLogger())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "community_posts" WHERE "community_posts"."id" = $1 ORDER BY "community_posts"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(newRows("id", "user_id", "org_id", "title", "content", "images", "post_type", "like_count", "comment_count", "status", "created_at").
			AddRow(1, 1, 0, "t", "c", "[]", "story", 0, 0, "published", time.Now()))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "post_comments"`)).WillReturnError(errors.New("db write failed"))
	mock.ExpectRollback()

	_, err := svc.Create(1, 1, "hi")
	if err == nil || !strings.Contains(err.Error(), "comment create") {
		t.Fatalf("expected preserved business error, got %v", err)
	}
}

func TestCommentDeletePreservesError(t *testing.T) {
	db, mock := newServiceDB(t)
	svc := NewCommentService(db, repository.NewPostCommentRepository(db), repository.NewCommunityPostRepository(db), newTestLogger())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "post_comments" WHERE "post_comments"."id" = $1 ORDER BY "post_comments"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(newRows("id", "post_id", "user_id", "content", "created_at").
			AddRow(1, 1, 1, "hi", time.Now()))
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "post_comments"`)).WillReturnError(errors.New("db delete failed"))

	if err := svc.Delete(1, 1); err == nil {
		t.Fatalf("expected preserved delete error, got nil")
	}
}
