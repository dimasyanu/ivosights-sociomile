<script setup lang="ts">
import { onBeforeMount, ref } from "vue"

const prop = defineProps({
  name: {
    type: String,
    required: true,
  },
  fill: {
    type: String,
    default: "transparent",
  },
  stroke: {
    type: String,
    default: "#333333",
  },
  size: {
    type: Number,
    default: 24,
  },
})

const icon = ref("")
onBeforeMount(
  async () =>
    (icon.value = (
      await import("../icons/" + prop.name + ".html?raw")
    ).default),
)
</script>
<template>
  <svg
    xmlns="http://www.w3.org/2000/svg"
    :fill="prop.fill"
    stroke-width="1.5"
    :width="prop.size"
    :height="prop.size"
    :stroke="prop.stroke"
    viewBox="0 0 24 24"
    v-html="icon"
  ></svg>
</template>
