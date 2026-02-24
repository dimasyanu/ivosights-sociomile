<script setup lang="ts">
import { type Menu } from "@/data/admin-menus"
import { lg, x2l, xl } from "@/data/media-queries"
import { useGlobalStore } from "@/stores/global"
import { useRouter } from "vue-router"
import SvgIcon from "./SvgIcon.vue"

const props = defineProps<{
  menu: Menu
  baseRoute: string
}>()
const globalStore = useGlobalStore()
const router = useRouter()
const goTo = (routeName: string) => {
  router.push({ name: routeName })
  if (lg.value || xl.value || x2l.value) return
  globalStore.isSidebarOpened = false
}
</script>
<template>
  <div
    v-if="!props.menu.children || props.menu.children.length < 1"
    :to="props.menu.name"
    @click="goTo(menu.name)"
  >
    <SvgIcon v-if="menu.icon" :name="menu.icon" :size="16" />
    <span>{{ props.menu.label }}</span>
  </div>
  <details v-else>
    <summary>
      <SvgIcon v-if="menu.icon" :name="menu.icon" :size="16" />
      <span>{{ props.menu.label }}</span>
    </summary>
    <ul>
      <li v-for="(m, i) in props.menu.children" :key="i">
        <SidebarMenuNode :menu="m" :baseRoute="`${baseRoute}.${m.name}`" />
      </li>
    </ul>
  </details>
</template>
