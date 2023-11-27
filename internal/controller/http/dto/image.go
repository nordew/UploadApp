package dto

type GetAllImageDTO struct {
	ID string `json:"id"`
}

type GetImageBySizeDTO struct {
	ID   string `json:"id"`
	Size int    `json:"size"`
}
