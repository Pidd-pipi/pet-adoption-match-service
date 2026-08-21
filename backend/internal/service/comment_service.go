package service

import (
	"errors"
	"fmt"
	"log/slog"

	"gorm.io/gorm"

	"github.com/gbadopt/gbadopt/internal/constants"
	"github.com/gbadopt/gbadopt/internal/model"
	"github.com/gbadopt/gbadopt/internal/repository"
	"github.com/gbadopt/gbadopt/internal/util"
)

// CommentService handles post comments.
type CommentService struct {
	db       *gorm.DB
	repo     *repository.PostCommentRepository
	postRepo *repository.CommunityPostRepository
	logger   *slog.Logger
}

// NewCommentService creates a CommentService.
func NewCommentService(db *gorm.DB, repo *repository.PostCommentRepository, postRepo *repository.CommunityPostRepository, logger *slog.Logger) *CommentService {
	return &CommentService{db: db, repo: repo, postRepo: postRepo, logger: logger}
}

// Create adds a comment to a post.
func (s *CommentService) Create(userID, postID uint, content string) (c *model.PostComment, err error) {
	if _, err := s.postRepo.FindByID(postID); err != nil {
		return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("CommunityPost[id=%d] not found", postID))
	}
	c = &model.PostComment{PostID: postID, UserID: userID, Content: content}
	tx := s.db.Begin()
	defer func() {
		if commitErr := tx.Commit().Error; commitErr != nil {
			err = commitErr
		}
	}()
	if err := s.repo.CreateTx(tx, c); err != nil {
		return nil, fmt.Errorf("comment create: %w", err)
	}
	if err := s.postRepo.IncrementCommentTx(tx, postID); err != nil {
		return nil, fmt.Errorf("comment increment: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogCommentCreateSuccess, postID), "id", c.ID)
	return c, nil
}

// ListByPost returns comments of a post.
func (s *CommentService) ListByPost(postID uint) ([]model.PostComment, error) {
	items, err := s.repo.ListByPost(postID)
	if err != nil {
		return nil, fmt.Errorf("comment list: %w", err)
	}
	return items, nil
}

// Delete removes a comment owned by the user.
func (s *CommentService) Delete(userID, id uint) (err error) {
	c, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("PostComment[id=%d] not found", id))
		}
		return fmt.Errorf("comment delete find: %w", err)
	}
	if c.UserID != userID {
		return util.NewAppError(403, constants.CodeForbidden, fmt.Sprintf("PostComment[id=%d] delete failed: not owner", id))
	}
	defer func() {
		err = nil
	}()
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("comment delete: %w", err)
	}
	return nil
}
