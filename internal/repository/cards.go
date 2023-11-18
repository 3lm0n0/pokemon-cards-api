package repository

import (
	"cards/internal/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CardsRepository interface {
	GetAllCards(ctx context.Context, ids []string, size int, page int, filters map[string]string) ([]models.Card, error)
	GetCards(ctx context.Context, ids []string, filters map[string]string) ([]models.Card, error)
	GetCard(ctx context.Context, id string) (models.Card, error)
	CreateCard(ctx context.Context, card models.Card) (models.Card, error)
	UpdateCard(ctx context.Context, card models.Card) error
	DeleteCard(ctx context.Context, id string) error
}

type CardsDB struct {
	db *gorm.DB
}

func NewCardsRepository(db *gorm.DB) CardsRepository {
	return &CardsDB{db: db}
}

func (orm *CardsDB) GetAllCards(ctx context.Context, ids []string, size int, page int, filters map[string]string) ([]models.Card, error) {
	var cards []models.Card
	result := orm.db.WithContext(ctx).Where(filters).Limit(size).Offset(page).Find(&cards)

	return cards, result.Error
}

func (orm *CardsDB) GetCards(ctx context.Context, ids []string, filters map[string]string) ([]models.Card, error) {
	var cards []models.Card
	result := orm.db.WithContext(ctx).Where("id IN ?", ids).Find(&cards)

	return cards, result.Error
}

func (orm *CardsDB) GetCard(ctx context.Context, id string) (models.Card, error) {
	var card models.Card
	result := orm.db.WithContext(ctx).First(&card, uuid.MustParse(id)) // string to uuid.

	return card, result.Error
}

func (orm *CardsDB) CreateCard(ctx context.Context, card models.Card) (models.Card, error) {
	result := orm.db.WithContext(ctx).Create(&card)

	return card, result.Error
}

func (orm *CardsDB) UpdateCard(ctx context.Context, card models.Card) error {
	return orm.db.WithContext(ctx).Model(&card).Updates(&card).Error
}

func (orm *CardsDB) DeleteCard(ctx context.Context, id string) error {
	return orm.db.WithContext(ctx).Where("id = ?", uuid.MustParse(id)).Delete(&models.Card{}).Error // string to uuid.
}