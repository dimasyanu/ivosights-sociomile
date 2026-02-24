<script setup lang="ts">
import { useGlobalStore } from "@/stores/global"
import { adminMenus } from "@/data/admin-menus"
import SidebarMenuNode from "./SidebarMenuNode.vue"
import { onMounted } from "vue"
import { md, xs } from "@/data/media-queries"
import SvgIcon from "./SvgIcon.vue"

const globalStore = useGlobalStore()
const basePath = import.meta.env.VITE_APP_ADMIN_PATH_PREFIX

const mouseEntered = () => {
  if (globalStore.isSidebarOpened || globalStore.isSidebarPinned) return
  globalStore.isSidebarOpened = true
}

const mouseLeaved = () => {
  if (globalStore.isSidebarPinned) return
  globalStore.isSidebarOpened = false
}
onMounted(() => {
  if (xs.value || md.value) {
    globalStore.isSidebarPinned = false
    globalStore.isSidebarOpened = false
  }
})
</script>

<template>
  <div id="sidebar" class="shadow-sm" :class="{
    hidden: !globalStore.isSidebarOpened,
    pinned: globalStore.isSidebarPinned,
  }" @mouseleave="mouseLeaved" @focusout="mouseLeaved">
    <div>
      <RouterLink to="/">
        <img src="/app.svg" alt="Logo" />
      </RouterLink>
      <span @click="globalStore.toggleSidebarPin">
        <SvgIcon name="pin" v-if="globalStore.isSidebarPinned" class="origin-center ml-1" fill="#7b7b7b" :width="18" />
        <SvgIcon name="unpin" v-else class="origin-center" fill="#7b7b7b" :width="18" />
      </span>
    </div>
    <div>
      <ul class="menu">
        <template v-for="(menu, i) in adminMenus" :key="i">
          <li v-if="menu.isReadOnly" class="menu-label">{{ menu.label }}</li>
          <li v-else>
            <SidebarMenuNode :menu="menu" :baseRoute="basePath" />
          </li>
        </template>
      </ul>
    </div>
    <div></div>
    <div class="hover-area" @mouseenter="mouseEntered" @click="mouseLeaved"></div>
  </div>
</template>

<style lang="sass" scoped>
@use '@/sass/size.scss' as size

#sidebar
  position: fixed
  display: flex
  flex-direction: column
  min-width: size.$sidebarWidth
  height: 100vh
  left: 0
  max-height: 100vh
  background-color: var(--color-base-200)
  border-right: 1px solid var(--color-base-100)
  transition: margin .2s ease-in-out, left .2s ease-in-out
  z-index: 10
  >div:first-child
    display: flex
    justify-content: space-between
    align-items: center
    padding-inline-start: calc(var(--spacing) * 5)
    padding-inline-end: calc(var(--spacing) * 5)
    min-height: calc(var(--spacing) * 16)
    >a
      height: 32px
      img
        height: 100%
    span
      display: none
      cursor: pointer
      max-width: 18px
      max-height: inherit
  >div:nth-child(2)
    padding: 0 calc(var(--spacing))
    .menu
      width: 100%
      font-family: var(--font-sans)
      font-weight: var(--font-weight-medium)
      .menu-label
        color: color-mix(in oklab, var(--color-base-content)70%, transparent)
        padding: 0 calc(var(--spacing)*3)
      li + .menu-label
        margin-top: calc(var(--spacing) * 2.5)
  >.hover-area
    display: block
    position: absolute
    left: 100%
    height: 100%
    width: calc(100vw - 100%)
    // background-color: rgb(50 50 50 / 45%)
    background-color: transparent
    transition: display .1s linear .25s
    -webkit-transition: display .1s linear .25s
  &.hidden
    >.hover-area
      display: none
      transition: display .1s linear .25s
      -webkit-transition: display .1s linear .25s
    &.pinned
      margin-left: -(size.$sidebarWidth)
    &:not(.pinned)
      left: -(size.$sidebarWidth)

@media screen and (min-width: size.$smMin)
  #sidebar
    >div:first-child
      span
        display: block
    &.pinned
      position: relative
      left: initial
    &.hidden
      >.hover-area
        display: block
    >.hover-area
      width: calc(var(--spacing) * 2)
      background-color: transparent
</style>
