package db

type ComputingBackend struct {
	BackendType string `json:"backend_type"`
	BackendInfo string `json:"backend_info"`
}

func NewComputingBackend(BackendType string, BackendInfo string) *ComputingBackend {
	backend := &ComputingBackend{
		BackendType: BackendType,
		BackendInfo: BackendInfo,
	}

	return backend
}
