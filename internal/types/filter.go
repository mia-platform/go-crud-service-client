package types

type Filter struct {
	FieldsQuery map[string]string
	MongoQuery  map[string]any
	Limit       int
	Projection  []string
	Skip        int
}
