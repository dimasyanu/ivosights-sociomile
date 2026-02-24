<script setup lang="ts">
import categoriesApi from "@/api/server/categories-api"
import conversation from "@/api/server/conversation"
import { toastHelper } from "@/api/toast-helper"
import FormModal from "@/components/FormModal.vue"
import PageHeading from "@/components/PageHeading.vue"
import TableData from "@/components/TableData.vue"
import type { Category } from "@/models/category"
import type { Item } from "@/models/item"
import type { TableHeader } from "@/models/table-header"
import {
  defineAsyncComponent,
  onMounted,
  ref,
  shallowRef,
  type Component,
} from "vue"

const items = ref<Item[]>([])
const headers: TableHeader[] = [
  { label: "Tenant", name: "tenant" },
  { label: "Title", name: "name" },
  { label: "Description", name: "description" },
  { label: "Created At", name: "createdAt" },
]
const total = ref<number>(0)
const pageSize = ref<number>(0)
const page = ref<number>(0)
const modal = ref<boolean>(false)
const modalBackdrop = ref<boolean>(false)
const modalTitle = ref<string>("Show")
const modalForm = shallowRef<Component | undefined>()
const modalData = ref<Category>({
  id: 0,
  name: "",
})

const loadItems = async () => {
  const result = await conversation.getConversations()
}

/// Create item
const create = () => {}

/// Show item detail
const show = async (id: number) => {
  modalData.value = await categoriesApi.getItem(id)
  if (!modalData.value)
    toastHelper.danger("Item not found. Please try reload the page.")

  modalTitle.value = "Category Detail"
  modalForm.value = defineAsyncComponent({
    loader: () => import("@/components/category/CategoryDetail.vue"),
  })
  modalBackdrop.value = true
  modal.value = true
}

/// Edit item
const edit = async (id: number) => {
  modalData.value = await categoriesApi.getItem(id)
  if (!modalData.value)
    toastHelper.danger("Item not found. Please try reload the page.")

  modalTitle.value = "Edit Category"
  modalBackdrop.value = false
  modalForm.value = defineAsyncComponent({
    loader: () => import("@/components/category/CategoryForm.vue"),
  })
  modal.value = true
}

/// Delete item
const remove = async (id: number) => {
  modalData.value = await categoriesApi.getItem(id)
  if (!modalData.value)
    toastHelper.danger("Item not found. Please try reload the page.")

  modalTitle.value = "Confirm Delete"
  modalForm.value = defineAsyncComponent({
    loader: () => import("@/components/category/CategoryDelete.vue"),
  })
  modalBackdrop.value = true
  modal.value = true
}

const closeAndRefresh = async () => {
  modal.value = false
  await loadItems()
}
onMounted(loadItems)
</script>
<template>
  <main>
    <PageHeading title="Categories" />
    <TableData
      :headers="headers"
      :data="items"
      :total-data="total"
      :page="page"
      :page-size="pageSize"
      @show-item="show"
      @create-item="create"
      @edit-item="edit"
      @delete-item="remove"
    ></TableData>
    <FormModal :title="modalTitle" :backdrop="modalBackdrop" v-model="modal">
      <component
        v-if="modalForm"
        :is="modalForm"
        :data="modalData"
        @requestRefresh="closeAndRefresh"
      ></component>
    </FormModal>
  </main>
</template>
