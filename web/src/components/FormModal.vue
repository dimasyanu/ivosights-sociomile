<script setup lang="ts">
import { ref, watch } from "vue"

const props = defineProps<{
  title: string
  backdrop?: boolean | undefined
}>()

const value = defineModel<boolean>()
const formModal = ref<HTMLDialogElement | undefined>()

const close = () => {
  value.value = false
}

const backdropClick = () => {
  if (!props.backdrop) return
  close()
}

watch(value, (val) => {
  if (!val) {
    formModal.value?.close()
    return
  }
  formModal.value?.showModal()
})
</script>
<template>
  <dialog ref="formModal" class="modal" @close="close" @cancel="close">
    <div class="modal-box">
      <form method="dialog">
        <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">
          ✕
        </button>
      </form>
      <h3 class="text-lg font-bold" v-html="title"></h3>
      <slot></slot>
    </div>
    <label class="modal-backdrop" for="form_modal" @click="backdropClick"
      >Close</label
    >
  </dialog>
</template>
<style lang="sass" scoped>
.modal
  .modal-box
    h3 + *
      padding-block: calc(var(--spacing) * 4)
</style>
