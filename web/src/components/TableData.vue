<script setup lang="ts">
import type { Item } from "@/models/item"
import type { TableHeader } from "@/models/table-header"
import SvgIcon from "./SvgIcon.vue"
import { useGlobalStore } from "@/stores/global"
import { computed, ref, watch } from "vue"

const props = defineProps<{
  headers: TableHeader[]
  data: Item[]
  totalData: number
  page: number
  pageSize: number
  pageSizeOptions?: number[] | undefined
}>()

const _pageSizeOptions = props.pageSizeOptions ?? [10, 25, 50, 100]
const globalStore = useGlobalStore()
const totalPage = computed(() => Math.ceil(props.totalData / props.pageSize))
const selectedPageSize = ref(_pageSizeOptions[0])
const searchKeyword = ref<string | undefined>()
let typingWatcher: number | undefined = undefined

const emit = defineEmits<{
  (e: "pageChange", page: number): void
  (e: "pageSizeChange", size: number): void
  (e: "search", keyword: string | undefined): void
  (e: "showItem", id: number): void
  (e: "createItem"): void
  (e: "editItem", id: number): void
  (e: "deleteItem", id: number): void
}>()

const typed = () => {
  if (typingWatcher) clearTimeout(typingWatcher)
  typingWatcher = setTimeout(() => {
    if (searchKeyword.value === "") searchKeyword.value = undefined
    emit("search", searchKeyword.value === "" ? undefined : searchKeyword.value)
    clearTimeout(typingWatcher)
  }, 750)
}
const changePage = (page: number) => {
  if (page === props.page) return
  emit("pageChange", page)
}
const pageSizeChange = (size: number) => {
  if (size === props.pageSize) return
  emit("pageSizeChange", size)
}

watch(selectedPageSize, (val) => {
  if (!val) return
  pageSizeChange(val)
})
</script>
<template>
  <div class="table-container card shadow-sm">
    <div class="table-filters">
      <div>
        <!-- Search input -->
        <label class="input input-sm">
          <input
            class="grow"
            type="text"
            name="global_search"
            id="global-search"
            placeholder="Search by keyword"
            v-model="searchKeyword"
            @keydown="typed"
          />
        </label>
      </div>
      <div></div>
    </div>
    <div class="table-data">
      <table class="table">
        <thead>
          <tr>
            <th class="text-center">No.</th>
            <template v-for="(header, i) in headers" :key="i">
              <th v-if="!header.hidden">{{ header.label }}</th>
            </template>
            <th class="text-center">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="data.length < 1">
            <td class="text-center" :colspan="headers.length + 2">
              No records
            </td>
          </tr>
          <template v-else v-for="(row, i) in data" :key="i">
            <tr>
              <td class="text-center">{{ i + 1 }}</td>
              <template v-for="(col, j) in row" :key="j">
                <td v-if="!headers[j]?.hidden">{{ col }}</td>
              </template>
              <td>
                <div class="text-center">
                  <span
                    class="btn btn-sm btn-ghost btn-circle"
                    @click="emit('showItem', row[0] as number)"
                  >
                    <SvgIcon
                      name="eye"
                      :size="16"
                      :stroke="globalStore.theme === 'light' ? '#333' : '#fff'"
                    />
                  </span>
                  <span
                    class="btn btn-sm btn-ghost btn-circle"
                    @click="emit('editItem', row[0] as number)"
                  >
                    <SvgIcon
                      name="edit"
                      :size="16"
                      :stroke="globalStore.theme === 'light' ? '#333' : '#fff'"
                    />
                  </span>
                  <span
                    class="btn btn-sm btn-ghost btn-circle"
                    @click="emit('deleteItem', row[0] as number)"
                  >
                    <SvgIcon name="trash" :size="16" stroke="#f31260" />
                  </span>
                </div>
              </td>
            </tr>
          </template>
        </tbody>
      </table>
    </div>
    <div class="table-footer">
      <div>
        <span>Per page</span>
        <select class="select select-sm" v-model="selectedPageSize">
          <option v-for="opt in _pageSizeOptions" :value="opt">
            {{ opt }}
          </option>
        </select>
      </div>
      <div>
        <div v-for="i = 1 in totalPage < 5 ? totalPage : 4" :key="i">
          <span
            class="btn btn-sm btn-ghost btn-circle"
            :class="{ active: i === page }"
            @click="changePage(i)"
            >{{ i }}</span
          >
        </div>
      </div>
      <span class="text-sm"
        >Showing <b>{{ page }}</b> to <b>{{ totalPage }}</b> of
        {{ totalData }} items</span
      >
    </div>
  </div>
</template>

<style lang="sass" scoped>
@use '@/sass/size' as size

.table-container
  background-color: var(--color-base-100)
  .table-filters
    display: flex
    padding-inline: calc(var(--spacing) * 4)
    padding-top: calc(var(--spacing) * 5)
    padding-bottom: calc(var(--spacing) * 5)
  .table-data
    overflow-x: auto
    width: 100%
    table.table
      td
        text-wrap: nowrap
      tr
        > *:first-child
          padding-inline: calc(var(--spacing) * 6)
        td:last-child span
          padding-inline: calc(var(--spacing) * 1.5)
  .table-footer
    display: flex
    padding-inline: calc(var(--spacing) * 4)
    padding-top: calc(var(--spacing) * 8)
    padding-bottom: 0
    justify-content: space-between
    align-items: center
    flex-wrap: wrap
    gap: 1rem
    >*
      margin-bottom: calc(var(--spacing) * 4)
    >div:first-child  /// Per page
      flex-basis: calc(var(--spacing) * 16)
      span
        display: none
        white-space: nowrap
      select
        text-align: end
    >div:nth-child(2) /// Page selection
      display: flex
      flex-grow: 1
      gap: calc(var(--spacing))
      justify-content: end
      .btn
        &.active
          color: var(--color-secondary-content)
          background-color: var(--color-secondary)
    >span /// Page & records information
      flex-grow: 1
      text-align: end
      opacity: 0.65
@media screen and (min-width: size.$smMin)
  .table-container
    .table-footer
      >div:first-child
        display: flex
        align-items: center
        gap: 1rem
        flex-basis: auto
        span
          display: inline-block
      >div:nth-child(2)
        justify-content: center
      >span
        flex-grow: 0
</style>
