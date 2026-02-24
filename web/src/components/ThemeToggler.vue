<script setup lang="ts">
import { useGlobalStore } from "@/stores/global"
import { computed } from "vue"
import SvgIcon from "./SvgIcon.vue"

defineProps<{
  size?: "sm" | "md" | "lg"
}>()

const globalStore = useGlobalStore()
const fgColor = computed(() =>
  globalStore.theme === "light" ? "#333" : "#fff",
)

const toggleTheme = () => {
  if (globalStore.theme !== "dark") {
    globalStore.changeTheme("dark")
    return
  }
  globalStore.changeTheme("light")
}
</script>
<template>
  <button
    class="btn btn-circle btn-ghost"
    :class="{ 'btn-sm': size === 'sm' }"
    aria-label="Tobble Theme"
    @click="toggleTheme"
  >
    <SvgIcon
      class="sun"
      name="sun"
      :class="{ active: globalStore.theme === 'light' }"
    />
    <SvgIcon
      class="moon"
      name="moon"
      :class="{ active: globalStore.theme === 'dark' }"
      :stroke="fgColor"
    />
  </button>
</template>

<style lang="sass" scoped>
button
  position: relative
  overflow: hidden
  box-sizing: border-box
  .sun,
  .moon
    position: absolute
    display: block
    left: 50%
    top: 50%
    min-width: 24px
    min-height: 24px
    opacity: 0.2
    transition: transform 0.25s ease-in-out, opacity 0.25s ease-in-out

    &.active
      opacity: 1
      transform: translate(-50%, -50%)
  .sun
    transform: translate(-50%, -100%)
    background-image: url("/icons/sun.svg")
  .moon
    transform: translate(-50%, 100%)
    background-image: url("/icons/moon.svg")
  &:hover
    cursor: pointer
</style>
