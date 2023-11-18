package service

import (
	models "cards/internal/models"
	"cards/internal/repository"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type Cards interface {
	GetCards(ctx context.Context, ids []string, p models.Pagination, filters map[string]string) ([]models.Card, error)
	CreateCard(ctx context.Context, request *http.Request) (*models.Card, error)
	UpdateCard(ctx context.Context, request *http.Request) (*models.Card, error)
	DeleteCard(ctx context.Context, id string) error
}

type CardsService struct {
	Repository repository.CardsRepository
	Pagination models.Pagination
}

func NewCardsService(cs CardsService) Cards {
	return &CardsService{
		Repository: cs.Repository,
		Pagination: models.Pagination{
			PageNumber: "1",
			PageSize: "10",
		},
	}
}

func (cs *CardsService) GetCards(ctx context.Context, ids []string, p models.Pagination, filters map[string]string) ([]models.Card, error) {
	// no ids requested.
	if len(ids) == int(1) && len(ids[0]) == int(0) {
		if p.PageNumber == "" {
			p.PageNumber = cs.Pagination.PageNumber

		}
		if p.PageSize == "" {
			p.PageSize = cs.Pagination.PageSize
		}		
		intPageNumber, _ := strconv.Atoi(p.PageNumber)
		intPageSize, _ := strconv.Atoi(p.PageSize)
		page := (intPageNumber - 1) * intPageSize

		cards, err := cs.Repository.GetAllCards(ctx, ids, intPageSize, page, filters)
		if err != nil {
			return nil, err
		}
		
		return cards, nil
	}
	// 1 id requested.
	if len(ids) == int(1) {
		var cards []models.Card
		card, err := cs.Repository.GetCard(ctx, ids[0])
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
		return cards, nil
	}
	// more than 1 id requested.
	cards, err := cs.Repository.GetCards(ctx, ids, filters)
	if err != nil {
		return nil, err
	}

	return cards, nil
}

func (cs *CardsService) CreateCard(ctx context.Context, request *http.Request) (*models.Card, error) {
	var card models.Card

	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(requestBody, &card)
	if err != nil {
		return nil, err
	}

	card.ID = uuid.New()

	card, err = cs.Repository.CreateCard(ctx, card)
	if err != nil {
		return nil, err
	}

	return &card, nil
}

func (cs *CardsService) UpdateCard(ctx context.Context, request *http.Request) (*models.Card, error) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	
	var card models.Card
	err = json.Unmarshal(requestBody, &card)
	if err != nil {
		return nil, err
	}

	err = cs.Repository.UpdateCard(ctx, card)
	if err != nil {
		return nil, err
	}

	newCard, err := cs.Repository.GetCard(ctx, card.ID.String())
	if err != nil {
		return nil, err
	}

	return &newCard, nil
}

func (cs *CardsService) DeleteCard(ctx context.Context, id string) error {
	return cs.Repository.DeleteCard(ctx, id)
}