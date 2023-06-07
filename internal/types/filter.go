package types

type Filter struct {
	MongoQuery map[string]any
	Limit      int
	Projection []string
	Skip       int
}
