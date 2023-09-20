<script setup lang="ts">
import { computed } from 'vue'
import * as echarts from 'echarts'
import type { EChartsOption } from 'echarts'
import { useClashTracingStore } from '@/stores'
import { formatTime } from '@/utils'
import BaseCard from './BaseCard.vue'
import BaseChart from './BaseChart.vue'

const store = useClashTracingStore()

const traffic = computed(() => store.traffic)

const formatTraffic = (byte: number) => {
  return (byte / 8 / 1024).toFixed(2)
}

const options = computed(() => {
  const opt: EChartsOption = {
    color: ['#8AB8FF', '#73BF69'],
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
        label: {
          backgroundColor: '#6a7985'
        }
      }
    },
    legend: {
      data: ['Up', 'Down']
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: [
      {
        type: 'category',
        boundaryGap: false,
        data: traffic.value.history.map((v) => formatTime(v.createTime)).reverse()
      }
    ],
    yAxis: [
      {
        type: 'value'
      }
    ],
    series: [
      {
        name: 'Up',
        type: 'line',
        stack: 'Total',
        smooth: true,
        lineStyle: {
          width: 0
        },
        showSymbol: false,
        areaStyle: {
          opacity: 0.8,
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            {
              offset: 0,
              color: 'rgb(128, 255, 165)'
            },
            {
              offset: 1,
              color: 'rgb(1, 191, 236)'
            }
          ])
        },
        emphasis: {
          focus: 'series'
        },
        data: traffic.value.history.map((v) => formatTraffic(v.up)).reverse()
      },
      {
        name: 'Down',
        type: 'line',
        stack: 'Total',
        smooth: true,
        lineStyle: {
          width: 0
        },
        showSymbol: false,
        areaStyle: {
          opacity: 0.8,
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            {
              offset: 0,
              color: 'rgb(0, 221, 255)'
            },
            {
              offset: 1,
              color: 'rgb(77, 119, 255)'
            }
          ])
        },
        emphasis: {
          focus: 'series'
        },
        data: traffic.value.history.map((v) => formatTraffic(v.down)).reverse()
      }
    ]
  }
  return opt
})
</script>

<template>
  <BaseCard title="实时流量">
    <BaseChart :options="options" />
  </BaseCard>
</template>

<style lang="less" scoped></style>
