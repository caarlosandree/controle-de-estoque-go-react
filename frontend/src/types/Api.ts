// Esta interface espelha a struct `domain.Metadata` do nosso backend.
export interface Metadata {
    total_records: number;
    current_page: number;
    page_size: number;
    total_pages: number;
}

// Esta interface genérica espelha a struct `domain.PaginatedResponse`.
export interface PaginatedResponse<T> {
    data: T[];
    metadata: Metadata;
}