package handler

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/gbadopt/gbadopt/internal/middleware"
	"github.com/gbadopt/gbadopt/internal/repository"
	"github.com/gbadopt/gbadopt/internal/service"
)

func newHandlerLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

func newHandlerDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

func TestDonationUsageHandlerStatus(t *testing.T) {
	db, mock := newHandlerDB(t)
	repo := repository.NewDonationRepository(db)
	usageRepo := repository.NewDonationUsageRepository(db)
	orgRepo := repository.NewOrganizationRepository(db)
	svc := service.NewDonationService(repo, usageRepo, orgRepo, newHandlerLogger())
	h := NewDonationHandler(svc, newHandlerLogger())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "donations" WHERE "donations"."id" = $1 ORDER BY "donations"."id" LIMIT $2`)).
		WithArgs(999, 1).WillReturnError(gorm.ErrRecordNotFound)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler(newHandlerLogger()))
	r.POST("/donations/usage", h.CreateUsage)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/donations/usage", strings.NewReader(`{"donation_id":999,"amount":50,"usage_desc":"x","proof_url":""}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
