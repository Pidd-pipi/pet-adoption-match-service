package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/gbadopt/gbadopt/internal/repository"
	"github.com/gbadopt/gbadopt/internal/service"
)

func newPetLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

func newPetDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("gorm open: %v", err)
	}
	return db, mock
}

func TestPetListHandlerContext(t *testing.T) {
	db, _ := newPetDB(t)
	repo := repository.NewPetRepository(db)
	svc := service.NewPetService(repo, nil, nil, newPetLogger())
	h := NewPetHandler(svc, newPetLogger())

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/pets?page=1&page_size=12", nil).WithContext(ctx)

	h.List(c)

	if len(c.Errors) == 0 {
		t.Fatalf("expected an error")
	}
	if !errors.Is(c.Errors.Last().Err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", c.Errors.Last().Err)
	}
}

func TestPetGetHandlerContext(t *testing.T) {
	db, _ := newPetDB(t)
	repo := repository.NewPetRepository(db)
	svc := service.NewPetService(repo, nil, nil, newPetLogger())
	h := NewPetHandler(svc, newPetLogger())

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/pets/1", nil).WithContext(ctx)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.Get(c)

	if len(c.Errors) == 0 {
		t.Fatalf("expected an error")
	}
	if !errors.Is(c.Errors.Last().Err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", c.Errors.Last().Err)
	}
}
