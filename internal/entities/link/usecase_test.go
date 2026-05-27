package link

import (
	"context"
	"errors"
	"sync"
	"testing"

	entryDomain "github.com/mkheyfets/ispro-app/internal/entities/entry/domain"
	"github.com/mkheyfets/ispro-app/internal/entities/link/domain"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/validator"
	"github.com/stretchr/testify/require"
)

type mockLinkRepository struct {
	mu    sync.Mutex
	items map[int64]*domain.Link
	seq   int64
	err   error
}

func newMockLinkRepository() *mockLinkRepository {
	return &mockLinkRepository{
		items: make(map[int64]*domain.Link),
	}
}

func (r *mockLinkRepository) setErr(err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.err = err
}

func (r *mockLinkRepository) List(_ context.Context) ([]*domain.Link, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.err != nil {
		return nil, r.err
	}
	result := make([]*domain.Link, 0, len(r.items))
	for _, l := range r.items {
		result = append(result, l)
	}
	return result, nil
}

func (r *mockLinkRepository) Get(_ context.Context, id int64) (*domain.Link, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.err != nil {
		return nil, r.err
	}
	l, ok := r.items[id]
	if !ok {
		return nil, errors.New("link not found")
	}
	return l, nil
}

func (r *mockLinkRepository) Create(_ context.Context, sourceID, targetID int64) (*domain.Link, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.err != nil {
		return nil, r.err
	}
	r.seq++
	l := &domain.Link{
		ID:       r.seq,
		SourceID: sourceID,
		TargetID: targetID,
	}
	r.items[l.ID] = l
	return l, nil
}

func (r *mockLinkRepository) Update(_ context.Context, id int64, sourceID, targetID int64) (*domain.Link, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.err != nil {
		return nil, r.err
	}
	l, ok := r.items[id]
	if !ok {
		return nil, errors.New("link not found")
	}
	l.SourceID = sourceID
	l.TargetID = targetID
	return l, nil
}

func (r *mockLinkRepository) Delete(_ context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.err != nil {
		return r.err
	}
	if _, ok := r.items[id]; !ok {
		return errors.New("link not found")
	}
	delete(r.items, id)
	return nil
}

type mockEntryRepo struct {
	mu    sync.Mutex
	items map[int64]*entryDomain.Entry
	seq   int64
	err   error
}

func newMockEntryRepo() *mockEntryRepo {
	return &mockEntryRepo{
		items: make(map[int64]*entryDomain.Entry),
	}
}

func (r *mockEntryRepo) List(_ context.Context) ([]*entryDomain.Entry, error) {
	return nil, nil
}

func (r *mockEntryRepo) Get(_ context.Context, id int64) (*entryDomain.Entry, error) {
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

func (r *mockEntryRepo) Create(_ context.Context, title, content string) (*entryDomain.Entry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	e := &entryDomain.Entry{ID: r.seq, Title: title, Content: content}
	r.items[e.ID] = e
	return e, nil
}

func (r *mockEntryRepo) Update(_ context.Context, id int64, title, content string) (*entryDomain.Entry, error) {
	return nil, nil
}

func (r *mockEntryRepo) Delete(_ context.Context, id int64) error {
	return nil
}

func TestCreateLink(t *testing.T) {
	tests := []struct {
		name      string
		cmd       CreateLinkCommand
		repoSetup func(repo *mockLinkRepository, entryRepo *mockEntryRepo)
		wantErr   bool
		wantNil   bool
		check     func(t *testing.T, result *domain.Link, err error)
	}{
		{
			name: "success",
			cmd:  CreateLinkCommand{SourceID: 1, TargetID: 2},
			repoSetup: func(repo *mockLinkRepository, entryRepo *mockEntryRepo) {
				_, _ = entryRepo.Create(context.Background(), "Entry", "Content")
				_, _ = entryRepo.Create(context.Background(), "Entry 2", "Content 2")
			},
			wantErr: false,
			wantNil: false,
			check: func(t *testing.T, result *domain.Link, _ error) {
				require.Equal(t, int64(1), result.SourceID)
				require.Equal(t, int64(2), result.TargetID)
			},
		},
		{
			name:    "validation_error_zero_source",
			cmd:     CreateLinkCommand{SourceID: 0, TargetID: 2},
			wantErr: true,
			wantNil: true,
			check: func(t *testing.T, _ *domain.Link, err error) {
				var vl *validation.ViolationList
				require.ErrorAs(t, err, &vl)
			},
		},
		{
			name:    "validation_error_zero_target",
			cmd:     CreateLinkCommand{SourceID: 1, TargetID: 0},
			wantErr: true,
			wantNil: true,
			check: func(t *testing.T, _ *domain.Link, err error) {
				var vl *validation.ViolationList
				require.ErrorAs(t, err, &vl)
			},
		},
		{
			name:    "equal_ids",
			cmd:     CreateLinkCommand{SourceID: 1, TargetID: 1},
			wantErr: true,
			wantNil: true,
			check: func(t *testing.T, _ *domain.Link, err error) {
				require.ErrorIs(t, err, validation.ErrIsEqual)
			},
		},
		{
			name: "entry_not_found",
			cmd:  CreateLinkCommand{SourceID: 1, TargetID: 2},
			repoSetup: func(repo *mockLinkRepository, entryRepo *mockEntryRepo) {
				_, _ = entryRepo.Create(context.Background(), "Entry", "Content")
			},
			wantErr: true,
			wantNil: true,
		},
		{
			name: "repo_error",
			cmd:  CreateLinkCommand{SourceID: 1, TargetID: 2},
			repoSetup: func(repo *mockLinkRepository, entryRepo *mockEntryRepo) {
				_, _ = entryRepo.Create(context.Background(), "Entry", "Content")
				_, _ = entryRepo.Create(context.Background(), "Entry 2", "Content 2")
				repo.setErr(errors.New("db error"))
			},
			wantErr: true,
			wantNil: true,
			check: func(t *testing.T, _ *domain.Link, err error) {
				require.Contains(t, err.Error(), "db error")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := newMockLinkRepository()
			entryRepo := newMockEntryRepo()
			if test.repoSetup != nil {
				test.repoSetup(repo, entryRepo)
			}
			uc := NewCreateLinkUseCase(repo, validator.Instance(), entryRepo)

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

func TestGetLink(t *testing.T) {
	tests := []struct {
		name      string
		query     GetLinkQuery
		repoSetup func(repo *mockLinkRepository)
		wantErr   bool
		wantNil   bool
		check     func(t *testing.T, result *domain.Link)
	}{
		{
			name:  "success",
			query: GetLinkQuery{ID: 1},
			repoSetup: func(repo *mockLinkRepository) {
				_, _ = repo.Create(context.Background(), 1, 2)
			},
			wantErr: false,
			wantNil: false,
			check: func(t *testing.T, result *domain.Link) {
				require.Equal(t, int64(1), result.SourceID)
				require.Equal(t, int64(2), result.TargetID)
			},
		},
		{
			name:    "not_found",
			query:   GetLinkQuery{ID: 999},
			wantErr: true,
			wantNil: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := newMockLinkRepository()
			if test.repoSetup != nil {
				test.repoSetup(repo)
			}
			uc := NewGetLinkUseCase(repo)

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

func TestListLinks(t *testing.T) {
	tests := []struct {
		name      string
		repoSetup func(repo *mockLinkRepository)
		wantLen   int
	}{
		{
			name: "success",
			repoSetup: func(repo *mockLinkRepository) {
				_, _ = repo.Create(context.Background(), 1, 2)
				_, _ = repo.Create(context.Background(), 3, 4)
			},
			wantLen: 2,
		},
		{
			name:    "empty",
			wantLen: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := newMockLinkRepository()
			if test.repoSetup != nil {
				test.repoSetup(repo)
			}
			uc := NewListLinksUseCase(repo)

			result, err := uc.Handle(context.Background())

			require.NoError(t, err)
			require.Len(t, result, test.wantLen)
		})
	}
}

func TestUpdateLink(t *testing.T) {
	tests := []struct {
		name      string
		cmd       UpdateLinkCommand
		repoSetup func(repo *mockLinkRepository, entryRepo *mockEntryRepo)
		wantErr   bool
		wantNil   bool
		check     func(t *testing.T, result *domain.Link, err error)
	}{
		{
			name: "success",
			cmd:  UpdateLinkCommand{ID: 1, SourceID: 3, TargetID: 4},
			repoSetup: func(repo *mockLinkRepository, entryRepo *mockEntryRepo) {
				_, _ = repo.Create(context.Background(), 1, 2)
				_, _ = entryRepo.Create(context.Background(), "Entry", "Content")
				_, _ = entryRepo.Create(context.Background(), "Entry 2", "Content 2")
				_, _ = entryRepo.Create(context.Background(), "Entry 3", "Content 3")
				_, _ = entryRepo.Create(context.Background(), "Entry 4", "Content 4")
			},
			wantErr: false,
			wantNil: false,
			check: func(t *testing.T, result *domain.Link, _ error) {
				require.Equal(t, int64(3), result.SourceID)
				require.Equal(t, int64(4), result.TargetID)
			},
		},
		{
			name:    "validation_error_zero_source",
			cmd:     UpdateLinkCommand{ID: 1, SourceID: 0, TargetID: 2},
			wantErr: true,
			wantNil: true,
			check: func(t *testing.T, _ *domain.Link, err error) {
				var vl *validation.ViolationList
				require.ErrorAs(t, err, &vl)
			},
		},
		{
			name:    "equal_ids",
			cmd:     UpdateLinkCommand{ID: 1, SourceID: 5, TargetID: 5},
			wantErr: true,
			wantNil: true,
			check: func(t *testing.T, _ *domain.Link, err error) {
				require.ErrorIs(t, err, validation.ErrIsEqual)
			},
		},
		{
			name: "entry_not_found",
			cmd:  UpdateLinkCommand{ID: 1, SourceID: 3, TargetID: 4},
			repoSetup: func(repo *mockLinkRepository, entryRepo *mockEntryRepo) {
				_, _ = repo.Create(context.Background(), 1, 2)
				_, _ = entryRepo.Create(context.Background(), "Entry", "Content")
			},
			wantErr: true,
			wantNil: true,
		},
		{
			name: "not_found",
			cmd:  UpdateLinkCommand{ID: 999, SourceID: 1, TargetID: 2},
			repoSetup: func(repo *mockLinkRepository, entryRepo *mockEntryRepo) {
				_, _ = entryRepo.Create(context.Background(), "Entry", "Content")
				_, _ = entryRepo.Create(context.Background(), "Entry 2", "Content 2")
			},
			wantErr: true,
			wantNil: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := newMockLinkRepository()
			entryRepo := newMockEntryRepo()
			if test.repoSetup != nil {
				test.repoSetup(repo, entryRepo)
			}
			uc := NewUpdateLinkUseCase(repo, validator.Instance(), entryRepo)

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

func TestDeleteLink(t *testing.T) {
	tests := []struct {
		name      string
		cmd       DeleteLinkCommand
		repoSetup func(repo *mockLinkRepository)
		wantErr   bool
		check     func(t *testing.T, repo *mockLinkRepository, err error)
	}{
		{
			name: "success",
			cmd:  DeleteLinkCommand{ID: 1},
			repoSetup: func(repo *mockLinkRepository) {
				_, _ = repo.Create(context.Background(), 1, 2)
			},
			wantErr: false,
			check: func(t *testing.T, repo *mockLinkRepository, _ error) {
				_, err := repo.Get(context.Background(), 1)
				require.Error(t, err)
			},
		},
		{
			name:    "not_found",
			cmd:     DeleteLinkCommand{ID: 999},
			wantErr: true,
			check: func(t *testing.T, _ *mockLinkRepository, err error) {
				require.Contains(t, err.Error(), "not found")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := newMockLinkRepository()
			if test.repoSetup != nil {
				test.repoSetup(repo)
			}
			uc := NewDeleteLinkUseCase(repo)

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
