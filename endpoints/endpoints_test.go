package endpoints

import (
	"io"
	"net/http"
	"testing"

	"github.com/LaughinKuma/tvdb-go-api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockClient is a mock implementation of the client.Client
type MockClient struct {
	mock.Mock
}

func (m *MockClient) Get(path string, result interface{}) error {
	args := m.Called(path, result)
	return args.Error(0)
}

func (m *MockClient) DoRequest(method, path string, body io.Reader) (*http.Response, error) {
	args := m.Called(method, path, body)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *MockClient) Post(path string, body interface{}, result interface{}) error {
    args := m.Called(path, body, result)
    return args.Error(0)
}

// Add the missing SetBaseURL method
func (m *MockClient) SetBaseURL(url string) {
	m.Called(url)
}

func TestGetSeriesByID(t *testing.T) {
	mockClient := new(MockClient)
	seriesID := 123
	expectedSeries := &models.Series{ID: seriesID, Name: "Test Series"}

	mockClient.On("Get", "/series/123", mock.AnythingOfType("*struct { Data models.Series }")).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).(*struct{ Data models.Series })
			arg.Data = *expectedSeries
		}).
		Return(nil)

	// Add expectation for SetBaseURL if it's used in the function
	mockClient.On("SetBaseURL", mock.Anything).Return()

	series, err := GetSeriesByID(mockClient, seriesID)

	assert.NoError(t, err)
	assert.Equal(t, expectedSeries, series)
	mockClient.AssertExpectations(t)
}

func TestGetSeriesSeasons(t *testing.T) {
	mockClient := new(MockClient)
	seriesID := 123
	expectedSeasons := []models.Season{{ID: 1, Name: "Season 1"}, {ID: 2, Name: "Season 2"}}

	mockClient.On("Get", "/series/123/seasons", mock.AnythingOfType("*struct { Data []models.Season }")).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).(*struct{ Data []models.Season })
			arg.Data = expectedSeasons
		}).
		Return(nil)

	// Add expectation for SetBaseURL if it's used in the function
	mockClient.On("SetBaseURL", mock.Anything).Return()

	seasons, err := GetSeriesSeasons(mockClient, seriesID)

	assert.NoError(t, err)
	assert.Equal(t, expectedSeasons, seasons)
	mockClient.AssertExpectations(t)
}

// Add other test functions similarly, ensuring to include the SetBaseURL expectation if necessary