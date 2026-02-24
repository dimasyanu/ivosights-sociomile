export type TableHeader = {
  name: string
  label: string
  sortable?: boolean | undefined
  hidden?: boolean | undefined
  align?: "center" | "left" | "right"
}
