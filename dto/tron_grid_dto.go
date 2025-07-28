package dto

type AccountResourceRequest struct {
	Address string `json:"address"`
}

type AccountResourceResponse struct {
	FreeNetUsed  int64 `json:"freeNetUsed"`
	FreeNetLimit int64 `json:"freeNetLimit"`
}
