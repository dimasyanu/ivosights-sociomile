<script setup lang="ts">
import { adminMenus, type Menu } from "@/data/admin-menus"
import type { Breadcrumb } from "@/models/breadcrumb"
import { onBeforeMount } from "vue"
import { useRoute } from "vue-router"
import SvgIcon from "./SvgIcon.vue"

defineProps<{
  title: string
}>()

const breadcrumbs: Breadcrumb[] = []
const route = useRoute()

onBeforeMount(() => {
  const routeSegments = (route.name! as string).split(".")
  if (routeSegments[0] === "admin") {
    let menus: Menu[] | undefined = adminMenus
    for (let i = 1; i < routeSegments.length; i++) {
      const segment = routeSegments[i]
      const menu: Menu | undefined = menus?.find((x) => x.name === segment)
      breadcrumbs.push({
        name:
          menu?.children && menu?.children.length > 0 ? undefined : menu?.alias,
        label: menu?.label ?? "",
        icon: menu?.icon,
      })
      menus = menu?.children
    }
    return
  }
})
</script>
<template>
  <div id="page-heading">
    <span>{{ title }}</span>
    <div class="breadcrumbs text-sm">
      <ul>
        <template v-for="item in breadcrumbs">
          <li v-if="!item.name" class="flex">
            <SvgIcon
              v-if="item.icon"
              :name="item.icon"
              :is="item.icon"
              :size="16"
            />
            {{ item.label }}
          </li>
          <li v-else>
            <RouterLink :to="{ name: item.name }">
              <SvgIcon v-if="item.icon" :name="item.icon" :size="16" />
              {{ item.label }}
            </RouterLink>
          </li>
        </template>
      </ul>
    </div>
  </div>
</template>
<style lang="sass" scoped>
@use '@/sass/size.scss' as size

#page-heading
  display: flex
  width: 100%
  align-items: center
  justify-content: space-between
  margin-bottom: calc(var(--spacing) * 6)
  > span
    font-size: var(--text-lg)
  ul
    li.flex
      gap: calc(var(--spacing) * 2)
      svg
        cursor: initial
  .breadcrumbs
    display: none

@media screen and (min-width: size.$smMin)
  #page-heading
    .breadcrumbs
      display: block
</style>
