import { defineStore } from "pinia"
import { v4 as uuid } from "uuid"

export type ToastType = "info" | "warning" | "danger"

export const useToastStore = defineStore("toast", {
  state: (): {
    items: { id: string; type: ToastType; message: string; timeout: number }[]
  } => ({
    items: [],
  }),
  actions: {
    push(type: ToastType, message: string, timeout: number = 2000) {
      // Add
      const item = {
        id: uuid(),
        type: type,
        message: message,
        timeout: timeout,
      }
      this.items.push(item)

      // Remove
      setTimeout(
        () => (this.items = this.items.filter((x) => x.id !== item.id)),
        timeout,
      )
    },
  },
})

export const toastHelper = {
  info: (message: string) => {
    const toastStore = useToastStore()
    toastStore.push("info", message)
  },
  danger: (message: string) => {
    const toastStore = useToastStore()
    toastStore.push("danger", message)
  },
  warning: (message: string) => {
    const toastStore = useToastStore()
    toastStore.push("warning", message)
  },
}
