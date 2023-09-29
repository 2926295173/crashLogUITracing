<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'

interface Props {
  minHeight?: string
  title: string
}

const defaultMinHeight = '120px'

const props = withDefaults(defineProps<Props>(), {
  minHeight: defaultMinHeight
})

const showSlot = ref(false)
const cardRef = ref<HTMLElement>()
const contentStyle = ref({
  height: 'auto',
  minHeight: props.minHeight
})

async function onResize() {
  showSlot.value = false
  nextTick(() => {
    const height = cardRef.value!.offsetHeight - 48
    contentStyle.value.height = height + 'px'
    showSlot.value = true
  })
}

watch(
  () => props.minHeight,
  (v) => {
    contentStyle.value.minHeight = v
    onResize()
  }
)

onMounted(() => {
  onResize()
  window.addEventListener('resize', onResize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', onResize)
})
</script>

<template>
  <div ref="cardRef" class="card">
    <div class="title">{{ title }}</div>
    <div :style="contentStyle" class="content">
      <slot v-if="showSlot" />
    </div>
  </div>
</template>

<style lang="less" scoped>
@import url(style.less);
</style>
