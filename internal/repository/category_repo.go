package repository

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/khai1301/moneywise-backend/internal/models"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *models.Category) error
	FindByUserID(userID string) ([]models.Category, error)
	FindByID(id string, userID string) (*models.Category, error)
	Update(category *models.Category) error
	Delete(id string, userID string) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) FindByUserID(userID string) ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.Where("user_id = ? OR is_system = ?", userID, true).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) FindByID(id string, userID string) (*models.Category, error) {
	var category models.Category
	result := r.db.Where("id = ? AND (user_id = ? OR is_system = ?)", id, userID, true).First(&category)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id string, userID string) error {
	result := r.db.Where("id = ? AND user_id = ? AND is_system = ?", id, userID, false).Delete(&models.Category{})
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == "23503" {
			return errors.New("Không thể xóa danh mục này vì vẫn đang còn giao dịch gắn với nó")
		}
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("Không thể xóa danh mục này (không tồn tại hoặc không đủ quyền)")
	}
	return nil
}
