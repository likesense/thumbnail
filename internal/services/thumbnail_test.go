package services

import (
	"testing"
	"thumbnail/internal/repositories"
	"thumbnail/internal/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockThumbnailRepository struct {
	mock.Mock
}

func (m *MockThumbnailRepository) Save(URL string, thumbnail []byte) error {
	args := m.Called(URL, thumbnail)
	return args.Error(0)
}

func (m *MockThumbnailRepository) Get(URL string) ([]byte, error) {
	args := m.Called(URL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockThumbnailRepository) Exists(URL string) (bool, error) {
	args := m.Called(URL)
	return args.Bool(0), args.Error(1)
}

func TestThumbnailService_GetThumbnail(t *testing.T) {
	mockRepo := new(MockThumbnailRepository)
	repos := &repositories.Repositories{Thumbnail: mockRepo}
	service := NewThumbnailService(repos)

	tests := []struct {
		name          string
		url           string
		setupMocks    func()
		expectedError bool
		expectedData  []byte
	}{
		{
			name: "Успешное получение существующего thumbnail",
			url:  "https://youtube.com/watch?v=abc123",
			setupMocks: func() {
				mockRepo.On("Exists", "https://youtube.com/watch?v=abc123").Return(true, nil)
				mockRepo.On("Get", "https://youtube.com/watch?v=abc123").Return([]byte("test data"), nil)
			},
			expectedError: false,
			expectedData:  []byte("test data"),
		},
		{
			name:          "Пустой URL",
			url:           "",
			setupMocks:    func() {},
			expectedError: true,
			expectedData:  nil,
		},
		{
			name: "Ошибка при проверке существования",
			url:  "https://youtube.com/watch?v=abc123",
			setupMocks: func() {
				mockRepo.On("Exists", "https://youtube.com/watch?v=abc123").Return(false, assert.AnError)
			},
			expectedError: true,
			expectedData:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.setupMocks()

			data, err := service.GetThumbnail(tt.url)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedData, data)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestThumbnailService_GetThumbnailAsync(t *testing.T) {
	mockRepo := new(MockThumbnailRepository)
	repos := &repositories.Repositories{Thumbnail: mockRepo}
	service := NewThumbnailService(repos)

	t.Run("Асинхронное получение нескольких thumbnails", func(t *testing.T) {
		urls := []string{
			"https://youtube.com/watch?v=abc123",
			"https://youtube.com/watch?v=def456",
		}

		mockRepo.On("Exists", mock.Anything).Return(true, nil)
		mockRepo.On("Get", mock.Anything).Return([]byte("test data"), nil)

		resultChan, err := service.GetThumbnailAsync(urls)
		assert.NoError(t, err)

		results := make([]utils.ThumbnailResult, 0)
		for result := range resultChan {
			results = append(results, result)
		}

		assert.Len(t, results, 2)
		for _, result := range results {
			assert.Empty(t, result.Error)
			assert.NotNil(t, result.Thumbnail)
		}
	})
}
