<script setup lang="ts">
import categoriesApi from "@/api/server/categories-api"
import type { Category } from "@/models/category"

const props = defineProps<{
  data: Category
}>()
const emit = defineEmits(["requestRefresh"])
const proceedDelete = () => {
  categoriesApi.deleteItem(props.data.id)
  emit("requestRefresh")
}
</script>
<template>
  <div class="delete-confirmation">
    <h2>Are you sure you want to delete this data?</h2>
    <div>
      <div>
        <span class="btn btn-sm btn-error" @click="proceedDelete">
          Yes, Delete
        </span>
      </div>
      <div>
        <form method="dialog">
          <button type="submit" class="btn btn-sm">Cancel</button>
        </form>
      </div>
    </div>
  </div>
</template>
<style lang="sass" scoped>
.delete-confirmation
  >div
    margin-top: calc(var(--spacing) * 6)
    display: flex
    justify-content: end
    gap: 1rem
    span.btn
      color: var(--color-base-100)
</style>
