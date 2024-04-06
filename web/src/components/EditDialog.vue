<script setup lang="ts">
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'

import {ref} from 'vue'

const open = ref(false)

document.addEventListener('open-edit-dialog',  _ => open.value = true)
document.addEventListener('close-edit-dialog',  _ => open.value = false)

</script>

<template>
  <Dialog :open="open">
    <DialogContent class="sm:max-w-[525px]">
      <DialogHeader>
        <DialogTitle class="text-4xl">Edit profile</DialogTitle>
        <DialogDescription class="text-xl">
          Make changes to your profile here.
        </DialogDescription>
      </DialogHeader>
      <div id="edit-container" />
    </DialogContent>
  </Dialog>
</template>
