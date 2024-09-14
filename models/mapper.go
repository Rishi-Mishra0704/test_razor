package models

// Mapper interface for mapping models to map[string]interface{}
type Mapper interface {
	ToMap() map[string]interface{}
}
