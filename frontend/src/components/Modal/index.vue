<script setup lang="ts">
interface Props {
  open: boolean
}

withDefaults(defineProps<Props>(), {})

const emits = defineEmits(['onCancel'])

const handleCancel = () => {
  emits('onCancel')
}
</script>

<template>
  <div v-if="open" class="modal">
    <div @click="handleCancel" class="mask"></div>
    <div class="content">
      <slot />
    </div>
  </div>
</template>

<style lang="less" scoped>
@import url(style.less);
</style>
