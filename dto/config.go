package dto

type CreateConfigDto struct {
	Content string `json:"content" binding:"required"`
}

type PatchConfigDto struct {
	CreateConfigDto
	ID uint `json:"id" binding:"required"`
}

type DeleteConfigDto struct {
	ID uint `uri:"id" binding:"required"`
}
