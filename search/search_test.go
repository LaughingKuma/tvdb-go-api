package search

import (
	"errors"
	"testing"

	"github.com/LaughinKuma/tvdb-go-api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockClient is a mock implementation of the ClientInterface
type MockClient struct {
	mock.Mock
}

func (m *MockClient) Get(path string, result interface{}) error {
	args := m.Called(path, result)
	return args.Error(0)
}

func TestSearch(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		mockResponse   []models.SearchResult
		mockError      error
		expectedResult []models.SearchResult
		expectedError  string
	}{
		{
			name:  "Successful search",
			query: "test series",
			mockResponse: []models.SearchResult{
				{ObjectID: "1", Type: "series", Name: "Test Series 1"},
				{ObjectID: "2", Type: "movie", Name: "Test Movie 1"},
			},
			mockError:      nil,
			expectedResult: []models.SearchResult{{ObjectID: "1", Type: "series", Name: "Test Series 1"}, {ObjectID: "2", Type: "movie", Name: "Test Movie 1"}},
			expectedError:  "",
		},
		{
			name:           "Empty result",
			query:          "nonexistent",
			mockResponse:   []models.SearchResult{},
			mockError:      nil,
			expectedResult: []models.SearchResult{},
			expectedError:  "",
		},
		{
			name:           "Client error",
			query:          "error test",
			mockResponse:   nil,
			mockError:      errors.New("client error"),
			expectedResult: nil,
			expectedError:  "search request failed: client error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockClient)

			mockClient.On("Get", mock.AnythingOfType("string"), mock.AnythingOfType("*struct { Data []models.SearchResult \"json:\\\"data\\\"\" }")).
				Run(func(args mock.Arguments) {
					result := args.Get(1).(*struct{ Data []models.SearchResult `json:"data"` })
					result.Data = tt.mockResponse
				}).
				Return(tt.mockError)

			results, err := Search(mockClient, tt.query)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedResult, results)

			mockClient.AssertExpectations(t)
		})
	}
}