package domain

// Metadata contém as informações de paginação.
type Metadata struct {
	TotalRecords int `json:"total_records"`
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	TotalPages   int `json:"total_pages"`
}

// PaginatedResponse é a estrutura genérica para respostas paginadas.
// Usamos `any` para que possamos reutilizá-la com qualquer tipo de dado (Produtos, Usuários, etc.).
type PaginatedResponse struct {
	Data     any      `json:"data"`
	Metadata Metadata `json:"metadata"`
}
