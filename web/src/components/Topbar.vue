<script setup lang="ts">
import { useGlobalStore } from "@/stores/global"
import ThemeToggler from "./ThemeToggler.vue"
import SvgIcon from "./SvgIcon.vue"

const globalStore = useGlobalStore()
</script>

<template>
  <div id="topbar" :class="{ expanded: !globalStore.isSidebarOpened }">
    <!-- Left side -->
    <div>
      <!-- Sidebar toggle button -->
      <button class="btn btn-square btn-ghost btn-sm" @click="globalStore.toggleSidebar">
        <SvgIcon name="arrow-left-to-line" :size="18" />
      </button>

      <!-- Search input -->
      <label class="input">
        <SvgIcon name="magnifier" stroke="#5b5b5b" />
        <input type="text" class="grow" name="global_search" id="global-search" placeholder="Search" />
      </label>
    </div>

    <!-- Right side -->
    <div>
      <!-- Theme button -->
      <ThemeToggler size="sm" />

      <!-- Notification button -->
      <button id="notification" class="btn btn-sm btn-ghost btn-circle p-1">
        <div class="indicator">
          <span class="indicator-item status status-error"></span>
          <SvgIcon name="bell" :stroke="globalStore.theme === 'dark' ? '#fff' : '#333'" />
        </div>
      </button>

      <!-- Profile button -->
      <div id="profile" class="btn btn-ghost">
        <div class="avatar">
          <div class="rounded-full">
            <img src="/images/user.png" alt="user" />
          </div>
        </div>
        <div>
          <p>Admin</p>
          <p>Administrator</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="sass" scoped>
@use "@/sass/size" as size
#topbar
  display: flex
  flex: 0 0 calc(var(--spacing) * 16)
  gap: 1rem
  align-items: center
  justify-content: space-between
  background-color: var(--color-base-200)
  padding-inline: calc(var(--spacing) * 3)
  width: 100%
  max-height: calc(var(--spacing) * l6)
  >div:first-child
    display: flex
    flex-direction: row
    align-items: center
    gap: calc(var(--spacing) * 3)
    button > svg
      transition: transform .25s ease-in-out
    label
      max-height: calc(var(--spacing) * 9)
      input
        &::placeholder
          font-weight: var(--font-weight-semibold)
  >div:last-child
    display: flex
    align-items: center
    gap: var(--spacing)
    #lang
      img
        max-width: initial
        height: calc(var(--spacing) * 5)
        width: calc(var(--spacing) * 5)
        border-radius: var(--radius-box)
        object-fit: cover
    #profile
      padding-inline: calc(var(--spacing)*1.5)
      text-align: start
      gap: calc(var(--spacing)*2)
      .avatar
        .rounded-full
          width: calc(var(--spacing) * 8)
      >div:last-child
        display: none
        p
          font-size: var(--text-sm)
          line-height: 1
        >p:last-child
          font-size: var(--text-xs)
          font-weight: var(--font-weight-normal)
          color: color-mix(in oklab,var(--color-base-content)50%,transparent)

  &.expanded
    >div:first-child
      button > svg
        transform: rotate(180deg)
  @media screen and (min-width: size.$smMin)
    >div:last-child
      #profile
        >div:last-child
          display: block
</style>
