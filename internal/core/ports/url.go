package ports

import (
	"URLShortner/internal/core/domain"
	"URLShortner/internal/infrastructure/persistence/ent"
)

func ToDomainUrl(url *ent.Urls) *domain.URL {
	return &domain.URL{
		ID:          url.ID,
		Destination: url.Destination,
		UserID:      url.UserID,
		ExpiresAt:   url.ExpiresAt,
		CreatedAt:   url.CreatedAt,
		UpdatedAt:   url.UpdatedAt,
		DeletedAt:   url.DeletedAt,
	}
}
