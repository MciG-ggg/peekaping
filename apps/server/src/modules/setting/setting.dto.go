package setting

type CreateDto struct {
}

type UpdateDto struct {
}

type CreateUpdateDto struct {
	Value string `json:"value"`
	Type  string `json:"type" validate:"required,oneof=string int bool json"`
}
