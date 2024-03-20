package crud

import "context"

type CrudClient[Resource any] interface {
	GetByID(ctx context.Context, id string, options Options) (*Resource, error)
	List(ctx context.Context, options Options) ([]Resource, error)
	Count(ctx context.Context, options Options) (int, error)
	Export(ctx context.Context, options Options) ([]Resource, error)
	PatchById(ctx context.Context, id string, body PatchBody, options Options) (*Resource, error)
	PatchMany(ctx context.Context, body PatchBody, options Options) (*Resource, error)
	PatchBulk(ctx context.Context, body PatchBulkBody, options Options) (int, error)
	Create(ctx context.Context, resource Resource, options Options) (string, error)
	CreateMany(ctx context.Context, resources []Resource, options Options) ([]CreatedResource, error)
	DeleteById(ctx context.Context, id string, options Options) error
	DeleteMany(ctx context.Context, options Options) (int, error)
	UpsertOne(ctx context.Context, body UpsertBody, options Options) (*Resource, error)
}
