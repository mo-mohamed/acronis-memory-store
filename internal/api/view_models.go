package api

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type SetRequest struct {
	Key        string      `json:"key"`
	Value      interface{} `json:"value"`
	TTLSeconds int         `json:"ttl_seconds"`
}

type UpdateRequest struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type PushRequest struct {
	Key  string      `json:"key"`
	Item interface{} `json:"item"`
}

type PopRequest struct {
	Key string `json:"key"`
}
