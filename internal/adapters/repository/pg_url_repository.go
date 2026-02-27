package repository

import (
	"URLShortner/internal/core/domain"
	"URLShortner/internal/core/ports"
	"URLShortner/internal/infrastructure/persistence/ent"
	"context"
	"sync"
)

type PgUrlRepository struct {
	mu     sync.RWMutex
	client *ent.Client
}

func NewPgUrlRepository(db *ent.Client) *PgUrlRepository {
	return &PgUrlRepository{
		client: db,
	}
}

func (r *PgUrlRepository) SaveUrl(ctx context.Context, url *domain.URL) error {
	_, err := r.client.Urls.
		Create().
		SetID(url.ID).
		SetDestination(url.Destination).
		SetUserID(url.UserID).
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *PgUrlRepository) UpdateUrl(ctx context.Context, url *domain.URL) error {
	// TODO
	return nil
}

func (r *PgUrlRepository) DeleteUrl(ctx context.Context, id string) error {
	// TODO
	return nil
}

func (r *PgUrlRepository) GetUrlById(ctx context.Context, id string) (*domain.URL, error) {
	url, err := r.client.Urls.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.ErrUrlNotFound
		}
		return nil, err
	}
	return ports.ToDomainUrl(url), nil
}

func (r *PgUrlRepository) ListUrlByUserId(ctx context.Context, userId int) ([]domain.URL, error) {
	// TODO
	return nil, nil
}
