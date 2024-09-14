package models

type Mapper interface {
	ToMap() map[string]interface{}
}
