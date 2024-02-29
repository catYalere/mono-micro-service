package repositories // import "github.com/catwashere/microservice/internal/database/repositories"

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/catwashere/microservice/internal/repositories"
	"io"
	"net/http"
	"net/url"
)

type Repository[T any] struct {
	BaseUrl string
	Path    string
}

func newGeneric[T any](_ context.Context, baseUrl string, path string) (repositories.IRepository[T], error) {
	return Repository[T]{
		BaseUrl: baseUrl,
		Path:    path,
	}, nil
}

// Get a list of resource
// The function is simply getting all entries in r.collection for the sake of example simplicity
func (r Repository[T]) Get(ctx context.Context, params map[string]interface{}) ([]T, error) {
	var result []T

	query := ""
	if len(params) > 0 {
		values := url.Values{}
		for k, v := range params {
			values.Add(k, fmt.Sprintf("%v", v))
		}
		query = "?" + values.Encode()
	}

	resp, err := http.Get(fmt.Sprintf("%s/%s%s", r.BaseUrl, r.Path, query))
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// GetOne resource based on its ID
func (r Repository[T]) GetOne(ctx context.Context, id string) (T, error) {
	var result T

	resp, err := http.Get(fmt.Sprintf("%s/%s/%s", r.BaseUrl, r.Path, id))
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// Create a new resource
func (r Repository[T]) Create(ctx context.Context, entity *T) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(entity)
	if err != nil {
		return err
	}

	resp, err := http.Post(fmt.Sprintf("%s/%s", r.BaseUrl, r.Path), "application/json", &buf)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &entity)
	if err != nil {
		return err
	}

	return nil
}

// Update a resource
func (r Repository[T]) Update(ctx context.Context, id string, entity *T) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(entity)
	if err != nil {
		return err
	}

	clientHttp := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s/%s", r.BaseUrl, r.Path, id), &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := clientHttp.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &entity)
	if err != nil {
		return err
	}

	return nil
}

// Delete a resource, soft delete by marking it as {"deleted": true}
func (r Repository[T]) Delete(ctx context.Context, id string) error {
	var buf bytes.Buffer
	clientHttp := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s/%s", r.BaseUrl, r.Path, id), &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := clientHttp.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
