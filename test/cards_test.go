package service_test

import (
	"context"
	"reflect"
	"testing"

	models "cards/internal/models"
	"cards/internal/repository"
	"cards/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func TestGetCards(t *testing.T) {
	mockRepo := &MockCardsRepository{} // Implement a mock repository to simulate database operations
	mockPagination := models.Pagination{PageNumber: "1", PageSize: "10"}
	mockFilters := map[string]string{}

	// Create a CardsService instance with the mock repository and pagination
	cardsService := service.CardsService{
		Repository: mockRepo,
		Pagination: mockPagination,
	}

	ctx := context.Background()

	t.Run("Get multiple cards by IDs", func(t *testing.T) {
		expectedIDs := []string{"0113d163-7f45-414d-961e-bb01db2e9eea", "5442df04-3f88-4c8c-8d36-60ebabc12f21", "f2957b26-4942-4ba5-aca5-1e37a3a7c5b7"}
		expectedCards := []models.Card{
			{ID: uuid.MustParse("0113d163-7f45-414d-961e-bb01db2e9eea"), Name: "Card 1"},
			{ID: uuid.MustParse("5442df04-3f88-4c8c-8d36-60ebabc12f21"), Name: "Card 2"},
			{ID: uuid.MustParse("f2957b26-4942-4ba5-aca5-1e37a3a7c5b7"), Name: "Card 3"},
		}

		// Mock the repository's GetCards method to return cards for provided IDs
		//mockRepo.CardsRepository.GetCards(ctx, expectedIDs, mockFilters)
		mockRepo.Mock.On("GetCards", ctx, expectedIDs, map[string]string{}).Return(expectedCards, nil)

		// Call the GetCards method with the provided IDs
		resultCards, err := cardsService.GetCards(ctx, expectedIDs, cardsService.Pagination, mockFilters)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Check if the returned cards match the expected cards
		if !reflect.DeepEqual(resultCards, expectedCards) {
			t.Errorf("Returned cards don't match expected cards")
		}
	})

}


// MockCardsRepository is a mock implementation of CardsRepository interface for testing purposes
type MockCardsRepository struct {
	repository.CardsRepository
	CalledMethods map[string]bool
	mock.Mock
}

// GetCards mocks the GetCards method of the repository interface
func (m *MockCardsRepository) GetCards(ctx context.Context, ids []string, filters map[string]string) ([]models.Card, error) {
	args := m.Mock.Called(ctx, ids, filters)
	return args.Get(0).([]models.Card), args.Error(1)
}
