package types

type Filter struct {
	Fields     map[string]string
	MongoQuery map[string]any
	Limit      int
	Projection []string
	Skip       int
}
