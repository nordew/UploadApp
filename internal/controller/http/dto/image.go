package dto

type GetAllImageDTO struct {
	ID string `json:"id"`
}

type GetImageBySizeDTO struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Size   int    `json:"size"`
}
