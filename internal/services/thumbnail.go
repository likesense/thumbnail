package services

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"
	"thumbnail/internal/repositories"
	"thumbnail/internal/utils"
)

type ThumbnailService struct {
	repo *repositories.Repositories
}

func NewThumbnailService(repo *repositories.Repositories) *ThumbnailService {
	return &ThumbnailService{
		repo: repo,
	}
}

func (ts *ThumbnailService) GetThumbnail(URL string) ([]byte, error) {
	if URL == "" {
		return nil, fmt.Errorf("empty url provided")
	}

	exists, err := ts.repo.Thumbnail.Exists(URL)
	if err != nil {
		return nil, fmt.Errorf("failed to check thumbnail exists: %w", err)
	}

	if exists {
		return ts.repo.Thumbnail.Get(URL)
	}

	thumbnailURL, err := ts.extractThumbnail(URL)
	if err != nil {
		return nil, fmt.Errorf("failed to extract thumbnail: %w", err)
	}
	thumbnail, err := ts.downloadThumbnail(thumbnailURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download thumbnail: %w", err)
	}
	if err = ts.repo.Thumbnail.Save(URL, thumbnail); err != nil {
		return nil, fmt.Errorf("failed to save thumbnail: %w", err)
	}

	return thumbnail, nil
}

func (ts *ThumbnailService) GetThumbnailAsync(URLs []string) (<-chan utils.ThumbnailResult, error) {
	result := make(chan utils.ThumbnailResult)

	semaphore := make(chan struct{}, 5)
	var wg sync.WaitGroup

	go func() {
		defer close(result)

		for _, url := range URLs {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()

				semaphore <- struct{}{}
				defer func() { <-semaphore }()
				exists, err := ts.repo.Thumbnail.Exists(url)
				if err != nil {
					result <- utils.ThumbnailResult{
						URL:   url,
						Error: fmt.Sprintf("failed to check thumbnail exists: %v", err),
					}
					return
				}
				if exists {
					data, err := ts.repo.Thumbnail.Get(url)
					if err != nil {
						result <- utils.ThumbnailResult{
							URL:   url,
							Error: fmt.Sprintf("failed to get thumbnail from cache: %v", err),
						}
						return
					}
					result <- utils.ThumbnailResult{
						URL:       url,
						Thumbnail: data,
					}
					return
				}
				data, err := ts.GetThumbnail(url)
				var errStr string
				if err != nil {
					errStr = err.Error()
				}

				result <- utils.ThumbnailResult{
					URL:       url,
					Thumbnail: data,
					Error:     errStr,
				}
			}(url)
		}

		wg.Wait()
	}()

	return result, nil
}

func (ts *ThumbnailService) extractThumbnail(URL string) (string, error) {
	re := regexp.MustCompile(`(?:youtube\.com\/watch\?v=|youtu.be\/)([^&]+)`)
	matches := re.FindStringSubmatch(URL)
	if len(matches) < 2 {
		return "", fmt.Errorf("invalid URL: %s", URL)
	}
	return fmt.Sprintf("https://img.youtube.com/vi/%s/maxresdefault.jpg", matches[1]), nil
}

func (ts *ThumbnailService) downloadThumbnail(URL string) ([]byte, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, fmt.Errorf("failed to download thumbnail: %w", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
