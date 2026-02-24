<script setup lang="ts">
import { RouterView, useRoute, useRouter } from "vue-router"
import Sidebar from "./components/Sidebar.vue"
import Topbar from "./components/Topbar.vue"
import { onBeforeMount, onMounted, ref } from "vue"
import { LS_THEME, useGlobalStore, type Theme } from "./stores/global"
import { useToastStore } from "./api/toast-helper"

const route = useRoute()
const router = useRouter()
const globalStore = useGlobalStore()
const toastStore = useToastStore()
const isAuthenticated = ref(false)

const checkAuth = () => {
  const token = localStorage.getItem("token")
  isAuthenticated.value = token !== null && token !== ""
  if (!isAuthenticated.value && route.name !== "login") {
    router.push({ name: "login" })
    return
  }
  if (isAuthenticated.value && route.name === "login") {
    router.push({ name: "conversations" })
  }
}
const checkTheme = () => {
  const storedTheme = localStorage.getItem(LS_THEME)
  globalStore.changeTheme(
    storedTheme === null ? "light" : (storedTheme as Theme),
  )
}

router.afterEach(checkAuth)
onBeforeMount(checkTheme)
onMounted(checkAuth)
</script>

<template>
  <div id="layout">
    <Sidebar v-if="isAuthenticated" />
    <div id="r">
      <Topbar v-if="isAuthenticated" />
      <RouterView />
    </div>
  </div>
  <div class="toast toast-end">
    <template v-for="toast in toastStore.items">
      <div
        class="alert"
        :class="{
          'alert-info': toast.type === 'info',
          'alert-warning': toast.type === 'warning',
          'alert-danger': toast.type === 'danger',
        }"
      >
        <span v-html="toast.message"></span>
      </div>
    </template>
  </div>
</template>

<style lang="sass" scoped>
#layout
  display: flex
  height: 100%
  #r
    display: flex
    max-width: 100vw
    flex-grow: 1
    flex-direction: column
</style>
