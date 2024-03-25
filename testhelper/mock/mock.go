package mock

import (
	"context"

	"github.com/mia-platform/go-crud-service-client"
)

type CRUD[Resource any] struct {
	GetByIDResult        *Resource
	GetByIDError         error
	GetByIDAssertionFunc func(ctx context.Context, id string, options crud.Options)

	ListResult        []Resource
	ListError         error
	ListAssertionFunc func(ctx context.Context, options crud.Options)

	CountResult        int
	CountError         error
	CountAssertionFunc func(ctx context.Context, options crud.Options)

	ExportResult        []Resource
	ExportError         error
	ExportAssertionFunc func(ctx context.Context, options crud.Options)

	PatchResult        *Resource
	PatchError         error
	PatchAssertionFunc func(ctx context.Context, id string, body crud.PatchBody, options crud.Options)

	PatchManyResult        int
	PatchManyError         error
	PatchManyAssertionFunc func(ctx context.Context, body crud.PatchBody, options crud.Options)

	PatchBulkResult        int
	PatchBulkError         error
	PatchBulkAssertionFunc func(ctx context.Context, body crud.PatchBulkBody, options crud.Options)

	CreateResult        string
	CreateError         error
	CreateAssertionFunc func(ctx context.Context, resource Resource, options crud.Options)

	CreateManyResult        []crud.CreatedResource
	CreateManyError         error
	CreateManyAssertionFunc func(ctx context.Context, resources []Resource, options crud.Options)

	DeleteByIDError         error
	DeleteByIDAssertionFunc func(ctx context.Context, id string, options crud.Options)

	DeleteManyResult        int
	DeleteManyError         error
	DeleteManyAssertionFunc func(ctx context.Context, options crud.Options)

	UpsertOneResult        *Resource
	UpsertOneError         error
	UpsertOneAssertionFunc func(ctx context.Context, body crud.UpsertBody, options crud.Options)
}

func (c *CRUD[Resource]) GetByID(ctx context.Context, id string, options crud.Options) (*Resource, error) {
	if c.GetByIDAssertionFunc != nil {
		c.GetByIDAssertionFunc(ctx, id, options)
	}
	return c.GetByIDResult, c.GetByIDError
}

func (c *CRUD[Resource]) List(ctx context.Context, options crud.Options) ([]Resource, error) {
	if c.ListAssertionFunc != nil {
		c.ListAssertionFunc(ctx, options)
	}
	return c.ListResult, c.ListError
}

func (c *CRUD[Resource]) Count(ctx context.Context, options crud.Options) (int, error) {
	if c.CountAssertionFunc != nil {
		c.CountAssertionFunc(ctx, options)
	}
	return c.CountResult, c.CountError
}

func (c *CRUD[Resource]) Export(ctx context.Context, options crud.Options) ([]Resource, error) {
	if c.ExportAssertionFunc != nil {
		c.ExportAssertionFunc(ctx, options)
	}
	return c.ExportResult, c.ExportError
}

func (c *CRUD[Resource]) PatchById(ctx context.Context, id string, body crud.PatchBody, options crud.Options) (*Resource, error) {
	if c.PatchAssertionFunc != nil {
		c.PatchAssertionFunc(ctx, id, body, options)
	}
	return c.PatchResult, c.PatchError
}

func (c *CRUD[Resource]) PatchMany(ctx context.Context, body crud.PatchBody, options crud.Options) (int, error) {
	if c.PatchManyAssertionFunc != nil {
		c.PatchManyAssertionFunc(ctx, body, options)
	}
	return c.PatchManyResult, c.PatchManyError
}

func (c *CRUD[Resource]) PatchBulk(ctx context.Context, body crud.PatchBulkBody, options crud.Options) (int, error) {
	if c.PatchBulkAssertionFunc != nil {
		c.PatchBulkAssertionFunc(ctx, body, options)
	}
	return c.PatchBulkResult, c.PatchBulkError

}

func (c *CRUD[Resource]) Create(ctx context.Context, resource Resource, options crud.Options) (string, error) {
	if c.CreateAssertionFunc != nil {
		c.CreateAssertionFunc(ctx, resource, options)
	}
	return c.CreateResult, c.CreateError

}

func (c *CRUD[Resource]) CreateMany(ctx context.Context, resources []Resource, options crud.Options) ([]crud.CreatedResource, error) {
	if c.CreateManyAssertionFunc != nil {
		c.CreateManyAssertionFunc(ctx, resources, options)
	}
	return c.CreateManyResult, c.CreateManyError
}

func (c *CRUD[Resource]) DeleteById(ctx context.Context, id string, options crud.Options) error {
	if c.DeleteByIDAssertionFunc != nil {
		c.DeleteByIDAssertionFunc(ctx, id, options)
	}
	return c.DeleteByIDError
}

func (c *CRUD[Resource]) DeleteMany(ctx context.Context, options crud.Options) (int, error) {
	if c.DeleteManyAssertionFunc != nil {
		c.DeleteManyAssertionFunc(ctx, options)
	}
	return c.DeleteManyResult, c.DeleteManyError
}

func (c *CRUD[Resource]) UpsertOne(ctx context.Context, body crud.UpsertBody, options crud.Options) (*Resource, error) {
	if c.UpsertOneAssertionFunc != nil {
		c.UpsertOneAssertionFunc(ctx, body, options)
	}
	return c.UpsertOneResult, c.UpsertOneError
}
