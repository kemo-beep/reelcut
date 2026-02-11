export interface ApiErrorDetail {
  field?: string
  message: string
}

export class ApiError extends Error {
  constructor(
    public code: string,
    message: string,
    public status: number,
    public details: ApiErrorDetail[] = []
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

export interface PaginationMeta {
  page: number
  per_page: number
  total_pages: number
  total_count: number
  has_next: boolean
  has_prev: boolean
}

export interface PaginatedResponse<T> {
  data: T
  pagination: PaginationMeta
}
