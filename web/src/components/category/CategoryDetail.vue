<script setup lang="ts">
import type { Category } from "@/models/category"

defineProps<{
  data: Category
}>()
</script>
<template>
  <div class="category-detail">
    <div>
      <span>Name</span>
      <span>: {{ data.name }}</span>
    </div>
    <div class="divider"></div>
    <div>
      <span>ID</span>
      <span>: {{ data.id }}</span>
      <span>Created By</span>
      <span>: {{ data.createdBy }}</span>
      <span>Created At</span>
      <span>: {{ data.createdAt }}</span>
      <span>Updated By</span>
      <span>: {{ data.updatedBy }}</span>
      <span>Updated At</span>
      <span>: {{ data.updatedAt }}</span>
    </div>
  </div>
</template>
<style lang="sass" scoped>
.category-detail
  padding-inline: calc(var(--spacing) * 2)
  >div
    display: grid
    grid-template-columns: 1fr 1fr

    &:last-child
      opacity: .5
  >.divider
    opacity: .4
    margin: var(--divider-m, .25rem 0)
    height: .5rem
</style>
