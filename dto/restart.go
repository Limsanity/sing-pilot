package dto

type RestartDto struct {
	ConfigId *uint `json:"configId" binding:"omitempty"`
}
