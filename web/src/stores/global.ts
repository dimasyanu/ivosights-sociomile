import { md, xs } from "@/data/media-queries"
import { defineStore } from "pinia"

export const LS_THEME: string = "theme"

export type Theme = "light" | "dark"

export const useGlobalStore = defineStore("globalStore", {
  state: () => ({
    theme: "light",
    isSidebarPinned: true,
    isSidebarOpened: true,
    accessToken: "",
  }),
  actions: {
    changeTheme(theme: Theme) {
      this.theme = theme
      this.applyThemeMode(theme)
    },
    applyThemeMode(t: Theme) {
      localStorage.setItem("theme", t)

      const html = document.getElementsByTagName("html")[0]
      if (!html || !html.dataset) return
      html.dataset.theme = t
    },
    toggleSidebar() {
      if (
        !xs.value &&
        !md.value &&
        !this.isSidebarPinned &&
        !this.isSidebarOpened
      )
        this.isSidebarPinned = true
      this.isSidebarOpened = !this.isSidebarOpened
    },
    toggleSidebarPin() {
      this.isSidebarPinned = !this.isSidebarPinned
    },
    setToken(token: string) {
      this.accessToken = token
      localStorage.setItem("token", token)
    },
  },
})
