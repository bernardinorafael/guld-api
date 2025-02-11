package org

import "context"

type RepositoryInterface interface {
	Insert(ctx context.Context, org Entity) error
	FindByID(ctx context.Context, id string) (*EntityWithOwner, error)
	FindBySlug(ctx context.Context, slug string) (*EntityWithOwner, error)
}

type ServiceInterface interface {
	CreateOrg(ctx context.Context, name, ownerId string) error
	GetOrgByID(ctx context.Context, orgId string) (*EntityWithOwner, error)
	GetOrgBySlug(ctx context.Context, slug string) (*EntityWithOwner, error)
}
