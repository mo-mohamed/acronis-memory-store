package api

type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

type SetRequest struct {
	Key        string `json:"key"`
	Value      any    `json:"value"`
	TTLSeconds int    `json:"ttl_seconds"`
}

type UpdateRequest struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type PushRequest struct {
	Key  string `json:"key"`
	Item any    `json:"item"`
}

type PopRequest struct {
	Key string `json:"key"`
}
