import type { TableHeader } from "./table-header"

export interface Paginated<T> {
  total: number
  perPage: number
  page: number
  headers: TableHeader[]
  items: T[]
}
