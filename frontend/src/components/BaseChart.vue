<script setup lang="ts">
import * as echarts from 'echarts'
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'

interface Props {
  width?: string
  height?: string
  options: any
}

const props = withDefaults(defineProps<Props>(), {
  width: '100%',
  height: '100%'
})

let chart: echarts.ECharts
const chartRef = ref<HTMLElement>()

function onResize() {
  chart && chart.resize()
}

watch(
  () => props.options,
  (opt) => {
    chart && chart.setOption(opt)
  }
)

onMounted(() => {
  window.addEventListener('resize', onResize)
  chart = echarts.init(chartRef.value, 'dark')
  chart.setOption({
    color: ['#8AB8FF', '#73BF69'],
    ...props.options
  })
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', onResize)
})
</script>

<template>
  <div :style="{ width, height }" ref="chartRef"></div>
</template>

<style lang="less" scoped></style>
