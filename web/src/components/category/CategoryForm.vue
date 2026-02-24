<script setup lang="ts">
import categoriesApi from "@/api/server/categories-api"
import type { Category } from "@/models/category"
import { ref } from "vue"

const props = defineProps<{
  data: Category
}>()

const name = ref(props.data.name)
const submit = () => {
  categoriesApi.saveItem({
    id: props.data.id,
    name: name.value,
  })
}
</script>
<template>
  <div>
    <form>
      <label for="id" class="input input-sm">ID</label>
      <input
        type="text"
        name="id"
        class="input input-sm text-right"
        :value="data.id"
        disabled
      />
      <label for="name" class="input input-sm">Name</label>
      <input type="text" name="name" class="input input-sm" v-model="name" />
    </form>
    <div class="actions">
      <div>
        <span class="btn btn-sm btn-success" @click="submit">Save</span>
      </div>
      <div>
        <form method="dialog">
          <button class="btn btn-sm" type="submit">Cancel</button>
        </form>
      </div>
    </div>
  </div>
</template>
<style lang="sass" scoped>
form
  display: grid
  grid-template-columns: 1fr 2fr
  row-gap: 1rem
  >label
    border-start-end-radius: 0
    border-end-end-radius: 0
    background-color: var(--color-base-300)
    border-inline-end: 0
    &+input
      border-start-start-radius: 0
      border-end-start-radius: 0
      border-color: var(--input-color)
      box-shadow: 0 1px color-mix(in oklab, var(--input-color) calc(var(--depth) * 10%), #0000) inset, 0 -1px oklch(100% 0 0 / calc(var(--depth) * 0.1)) inset
.actions
  display: flex
  justify-content: end
  gap: .5rem
  margin-top: calc(var(--spacing) * 6)
  form
    display: block
</style>
