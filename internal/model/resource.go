package model

// Resources ...
type Resources struct {
	Memory Resource
	CPU    Resource
}

// Resource ...
type Resource struct {
	Limit   string
	Request string
}
