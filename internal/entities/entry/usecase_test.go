package entry

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/mkheyfets/ispro-app/internal/entities/entry/domain"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/validator"
	"github.com/stretchr/testify/require"
)

type mockEntryRepository struct {
	mu    sync.Mutex
	items map[int64]*domain.Entry
	seq   int64
	err   error
}

func newMockEntryRepository() *mockEntryRepository {
	return &mockEntryRepository{
		items: make(map[int64]*domain.Entry),
	}
}

func (r *mockEntryRepository) setErr(err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.err = err
}

func (r *mockEntryRepository) List(_ context.Context) ([]*domain.Entry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.err != nil {
		return nil, r.err
	}
	result := make([]*domain.Entry, 0, len(r.items))
	for _, e := range r.items {
		result = append(result, e)
	}
	return result, nil
}

func (r *mockEntryRepository) Get(_ context.Context, id int64) (*domain.Entry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.err != nil {
		return nil, r.err
	}
	e, ok := r.items[id]
	if !ok {
		return nil, errors.New("entry not found")
	}
	return e, nil
}

func (r *mockEntryRepository) Create(_ context.Context, title, content string) (*domain.Entry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.err != nil {
		return nil, r.err
	}
	r.seq++
	e := &domain.Entry{
		ID:      r.seq,
		Title:   title,
		Content: content,
	}
	r.items[e.ID] = e
	return e, nil
}

func (r *mockEntryRepository) Update(_ context.Context, id int64, title, content string) (*domain.Entry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.err != nil {
		return nil, r.err
	}
	e, ok := r.items[id]
	if !ok {
		return nil, errors.New("entry not found")
	}
	e.Title = title
	e.Content = content
	return e, nil
}

func (r *mockEntryRepository) Delete(_ context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.err != nil {
		return r.err
	}
	if _, ok := r.items[id]; !ok {
		return errors.New("entry not found")
	}
	delete(r.items, id)
	return nil
}

func TestCreateEntry(t *testing.T) {
	tests := []struct {
		name      string
		cmd       CreateEntryCommand
		repoSetup func(repo *mockEntryRepository)
		wantErr   bool
		wantNil   bool
		check     func(t *testing.T, result *domain.Entry, err error)
	}{
		{
			name:    "success",
			cmd:     CreateEntryCommand{Title: "Title", Content: "Content"},
			wantErr: false,
			wantNil: false,
			check: func(t *testing.T, result *domain.Entry, _ error) {
				require.Equal(t, int64(1), result.ID)
				require.Equal(t, "Title", result.Title)
				require.Equal(t, "Content", result.Content)
			},
		},
		{
			name:    "validation_error",
			cmd:     CreateEntryCommand{Title: "", Content: "Content"},
			wantErr: true,
			wantNil: true,
			check: func(t *testing.T, _ *domain.Entry, err error) {
				var violationList *validation.ViolationList
				require.ErrorAs(t, err, &violationList)
			},
		},
		{
			name:    "repo_error",
			cmd:     CreateEntryCommand{Title: "Title", Content: "Content"},
			wantErr: true,
			wantNil: true,
			repoSetup: func(repo *mockEntryRepository) {
				repo.setErr(errors.New("db error"))
			},
			check: func(t *testing.T, _ *domain.Entry, err error) {
				require.Contains(t, err.Error(), "db error")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := newMockEntryRepository()
			if test.repoSetup != nil {
				test.repoSetup(repo)
			}
			uc := NewCreateEntryUseCase(repo, validator.Instance())

			result, err := uc.Handle(context.Background(), test.cmd)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if test.wantNil {
				require.Nil(t, result)
			} else {
				require.NotNil(t, result)
			}
			if test.check != nil {
				test.check(t, result, err)
			}
		})
	}
}

func TestGetEntry(t *testing.T) {
	tests := []struct {
		name      string
		query     GetEntryQuery
		repoSetup func(repo *mockEntryRepository)
		wantErr   bool
		wantNil   bool
		check     func(t *testing.T, result *domain.Entry)
	}{
		{
			name:  "success",
			query: GetEntryQuery{ID: 1},
			repoSetup: func(repo *mockEntryRepository) {
				_, _ = repo.Create(context.Background(), "Title", "Content")
			},
			wantErr: false,
			wantNil: false,
			check: func(t *testing.T, result *domain.Entry) {
				require.Equal(t, "Title", result.Title)
			},
		},
		{
			name:    "not_found",
			query:   GetEntryQuery{ID: 999},
			wantErr: true,
			wantNil: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := newMockEntryRepository()
			if test.repoSetup != nil {
				test.repoSetup(repo)
			}
			uc := NewGetEntryUseCase(repo)

			result, err := uc.Handle(context.Background(), test.query)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if test.wantNil {
				require.Nil(t, result)
			} else {
				require.NotNil(t, result)
			}
			if test.check != nil {
				test.check(t, result)
			}
		})
	}
}

func TestListEntries(t *testing.T) {
	tests := []struct {
		name      string
		repoSetup func(repo *mockEntryRepository)
		wantLen   int
	}{
		{
			name:    "success",
			wantLen: 2,
			repoSetup: func(repo *mockEntryRepository) {
				_, _ = repo.Create(context.Background(), "Title 1", "Content 1")
				_, _ = repo.Create(context.Background(), "Title 2", "Content 2")
			},
		},
		{
			name:    "empty",
			wantLen: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := newMockEntryRepository()
			if test.repoSetup != nil {
				test.repoSetup(repo)
			}
			uc := NewListEntriesUseCase(repo)

			result, err := uc.Handle(context.Background())

			require.NoError(t, err)
			require.Len(t, result, test.wantLen)
		})
	}
}

func TestUpdateEntry(t *testing.T) {
	tests := []struct {
		name      string
		cmd       UpdateEntryCommand
		repoSetup func(repo *mockEntryRepository)
		wantErr   bool
		wantNil   bool
		check     func(t *testing.T, result *domain.Entry, err error)
	}{
		{
			name: "success",
			cmd:  UpdateEntryCommand{ID: 1, Title: "New Title", Content: "New Content"},
			repoSetup: func(repo *mockEntryRepository) {
				_, _ = repo.Create(context.Background(), "Old Title", "Old Content")
			},
			wantErr: false,
			wantNil: false,
			check: func(t *testing.T, result *domain.Entry, _ error) {
				require.Equal(t, "New Title", result.Title)
				require.Equal(t, "New Content", result.Content)
			},
		},
		{
			name: "validation_error",
			cmd:  UpdateEntryCommand{ID: 1, Title: "", Content: "Content"},
			repoSetup: func(repo *mockEntryRepository) {
				_, _ = repo.Create(context.Background(), "Title", "Content")
			},
			wantErr: true,
			wantNil: true,
			check: func(t *testing.T, _ *domain.Entry, err error) {
				var violationList *validation.ViolationList
				require.ErrorAs(t, err, &violationList)
			},
		},
		{
			name:    "not_found",
			cmd:     UpdateEntryCommand{ID: 999, Title: "Title", Content: "Content"},
			wantErr: true,
			wantNil: true,
			check: func(t *testing.T, _ *domain.Entry, err error) {
				require.Contains(t, err.Error(), "not found")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := newMockEntryRepository()
			if test.repoSetup != nil {
				test.repoSetup(repo)
			}
			uc := NewUpdateEntryUseCase(repo, validator.Instance())

			result, err := uc.Handle(context.Background(), test.cmd)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if test.wantNil {
				require.Nil(t, result)
			} else {
				require.NotNil(t, result)
			}
			if test.check != nil {
				test.check(t, result, err)
			}
		})
	}
}

func TestDeleteEntry(t *testing.T) {
	tests := []struct {
		name      string
		cmd       DeleteEntryCommand
		repoSetup func(repo *mockEntryRepository)
		wantErr   bool
		check     func(t *testing.T, repo *mockEntryRepository, err error)
	}{
		{
			name: "success",
			cmd:  DeleteEntryCommand{ID: 1},
			repoSetup: func(repo *mockEntryRepository) {
				_, _ = repo.Create(context.Background(), "Title", "Content")
			},
			wantErr: false,
			check: func(t *testing.T, repo *mockEntryRepository, _ error) {
				_, err := repo.Get(context.Background(), 1)
				require.Error(t, err)
			},
		},
		{
			name:    "not_found",
			cmd:     DeleteEntryCommand{ID: 999},
			wantErr: true,
			check: func(t *testing.T, _ *mockEntryRepository, err error) {
				require.Contains(t, err.Error(), "not found")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := newMockEntryRepository()
			if test.repoSetup != nil {
				test.repoSetup(repo)
			}
			uc := NewDeleteEntryUseCase(repo)

			err := uc.Handle(context.Background(), test.cmd)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if test.check != nil {
				test.check(t, repo, err)
			}
		})
	}
}
