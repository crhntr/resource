package example

import "github.com/crhntr/resource"

type Resource struct {
	URI    string `json:"uri"`
	Branch string `json:"branch"`
}

type GetParams struct {
	IncludeZip bool `json:"include_zip"`
}

type PutParams struct {
	EnsureChecksum bool `json:"ensure_checksum"`
}

type Version struct {
	Ref string `json:"ref"`
}

var (
	_ resource.Get[Resource, GetParams, Version]
	_ resource.Put[Resource, PutParams, Version]
	_ resource.Check[Resource, Version]
)
