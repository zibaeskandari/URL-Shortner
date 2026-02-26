package domain

import (
	"context"
	"errors"
)

var (
	ErrUrlAlreadyExists        = errors.New("url already exists")
	ErrUrlNotFound             = errors.New("url not found")
	ErrUrlDoesNotBelongToUser  = errors.New("url does not belong to user")
	ErrUrlAlreadyDeleted       = errors.New("url already deleted")
	ErrUrlExpired              = errors.New("url expired")
	ErrUrlDestCannotBeModified = errors.New("url destination cannot be modified")
)

type UrlRepository interface {
	SaveUrl(ctx context.Context, url *URL) error
	UpdateUrl(ctx context.Context, url *URL) error
	DeleteUrl(ctx context.Context, id string) error
	GetUrlById(ctx context.Context, id string) (*URL, error)
	ListUrlByUserId(ctx context.Context, userId int) ([]URL, error)
}
