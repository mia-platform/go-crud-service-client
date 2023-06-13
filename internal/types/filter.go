package types

type Filter struct {
	FiledsQuery map[string]string
	MongoQuery  map[string]any
	Limit       int
	Projection  []string
	Skip        int
}
