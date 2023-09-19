<script setup lang="ts">
import { computed } from 'vue'
import type { EChartsOption } from 'echarts'
import { useClashTracingStore } from '@/stores'
import BaseCard from './BaseCard.vue'
import BaseChart from './BaseChart.vue'

const store = useClashTracingStore()

const proxy = computed(() =>
  store.proxyDial.proxyCounting.map((v) => ({ ...v, duration: formatNumber(v.duration) }))
)

const formatNumber = (n: number) => Number((n / 1000).toFixed(2))

const options = computed(() => {
  const opt: EChartsOption = {
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    legend: {
      data: ['duration', 'count']
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '0%',
      containLabel: true
    },
    xAxis: [
      {
        type: 'category',
        boundaryGap: false,
        axisLabel: {
          interval: 0,
          rotate: 25
        },
        data: proxy.value.map((v) => v.proxy)
      }
    ],
    yAxis: [
      {
        type: 'value'
      }
    ],
    series: [
      {
        name: 'duration',
        type: 'bar',
        barGap: 0,
        emphasis: {
          focus: 'series'
        },
        data: proxy.value.map((v) => v.duration)
      },
      {
        name: 'count',
        type: 'bar',
        emphasis: {
          focus: 'series'
        },
        data: proxy.value.map((v) => v.count)
      }
    ]
  }
  return opt
})
</script>

<template>
  <BaseCard title="代理使用次数">
    <BaseChart height="240px" :options="options" />
  </BaseCard>
</template>

<style lang="less" scoped></style>
