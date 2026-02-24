import type { Category } from "@/models/category"
import type { Item } from "@/models/item"
import type { Paginated } from "@/models/paginated"
import type { TableFilter } from "@/models/table-filter"
import type { TableHeader } from "@/models/table-header"

const LOCAL_CATEGORIES = "categories"
const headers: TableHeader[] = [
  {
    name: "id",
    label: "ID",
    hidden: true,
  },
  {
    name: "name",
    label: "Name",
  },
  {
    name: "createdAt",
    label: "Created At",
  },
  {
    name: "createdBy",
    label: "Created By",
  },
  {
    name: "updatedAt",
    label: "Updated At",
  },
  {
    name: "updatedBy",
    label: "Updated By",
  },
]
const initialItem: Item = [
  1,
  "Category 1",
  "2026-02-10 11:00",
  "admin",
  "2026-02-10 11:00",
  "admin",
]

const setInitialItems = (): void =>
  localStorage.setItem(LOCAL_CATEGORIES, JSON.stringify([initialItem]))

const getAllItems = (): Item[] => {
  const itemsStr = localStorage.getItem(LOCAL_CATEGORIES)
  let items: Item[] = []
  if (itemsStr === null) setInitialItems()
  else {
    items = JSON.parse(itemsStr)
    if (items.length < 1) setInitialItems()
  }
  return items
}

const defaultPageSize = 25
const defaultPage = 1

export default {
  getItems: (filter?: TableFilter | undefined): Promise<Paginated<Item>> =>
    new Promise((resolve, reject) => {
      try {
        const page = filter?.page ?? defaultPage
        const pageSize = filter?.pageSize ?? defaultPageSize
        const items = getAllItems().slice((page - 1) * pageSize, pageSize)

        resolve({
          items: items,
          headers: headers,
          page: 1,
          perPage: 25,
          total: items.length,
        })
      } catch (error) {
        reject(error)
      }
    }),
  saveItem: (model: Category): Promise<number> =>
    new Promise((resolve, reject) => {
      // const item: Item = [model.id, model.name]
      try {
        const items = getAllItems()
        if (model.id) {
          let item = items.find((x) => x[0] === model.id)
          if (!item) throw new Error("Item not found.")
          item[1] = model.name
          localStorage.setItem(LOCAL_CATEGORIES, JSON.stringify(items))
          resolve(model.id)
          return
        }

        const newId = Math.max(...items.map((x) => x[0] as number)) + 1
        const d = new Date()
        const dStr = `${d.getFullYear()}-${d.getMonth()}-${d.getDate()} ${d.getHours()}:${d.getMinutes()}`
        let item: Item = [newId, model.name, dStr, "admin", dStr, "admin"]
        items.push(item)
        localStorage.setItem(LOCAL_CATEGORIES, JSON.stringify(items))
        resolve(newId)
      } catch (error) {
        reject(error)
      }
    }),
  getItem: (id: number): Promise<Category> =>
    new Promise((resolve, reject) => {
      try {
        const items = getAllItems()
        const item = items.find((x) => x[0] === id)
        if (item === undefined) throw new Error("Item not found")
        resolve({
          id: item[0] as number,
          name: item[1] as string,
          createdAt: item[2] as string,
          createdBy: item[3] as string,
          updatedAt: item[4] as string,
          updatedBy: item[5] as string,
        })
      } catch (error) {
        reject(error)
      }
    }),
  deleteItem: (id: number): Promise<boolean> =>
    new Promise((resolve, reject) => {
      try {
        let items = getAllItems()
        items = items.filter((x) => x[0] !== id)
        localStorage.setItem(LOCAL_CATEGORIES, JSON.stringify(items))
        resolve(true)
      } catch (error) {
        reject(error)
      }
    }),
}
