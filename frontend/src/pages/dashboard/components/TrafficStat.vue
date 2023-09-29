<script setup lang="ts">
import { computed } from 'vue'
import * as echarts from 'echarts'
import type { EChartsOption } from 'echarts'
import { useClashTracingStore } from '@/stores'
import { formatTime, formatTraffic } from '@/utils'
import Card from '@/components/Card/index.vue'
import Chart from '@/components/Chart/index.vue'

const store = useClashTracingStore()

const traffic = computed(() => store.traffic)

const options = computed(() => {
  const opt: EChartsOption = {
    color: ['#8AB8FF', '#73BF69'],
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        const time = params[0].name
        const up = params[0].seriesName + '：' + formatTraffic(params[0].value)
        const down = params[1].seriesName + '：' + formatTraffic(params[1].value)
        const upColor = params[0].color
        const downColor = params[1].color
        const dom = /*html */ `
          <div>
            <div>${time}</div>
            <div><span style="width: 10px; height: 10px; margin-right: 4px; border-radius: 10px;display: inline-block; background-color: ${upColor}"></span>${up}</div>
            <div><span style="width: 10px; height: 10px; margin-right: 4px; border-radius: 10px;display: inline-block; background-color: ${downColor}"></span>${down}</div>
          </div>
        `
        return dom
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
        type: 'value',
        axisLabel: {
          formatter(value) {
            return formatTraffic(value)
          }
        }
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
        data: traffic.value.history.map((v) => v.up).reverse()
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
        data: traffic.value.history.map((v) => v.down).reverse()
      }
    ]
  }
  return opt
})
</script>

<template>
  <Card title="实时流量">
    <Chart :options="options" />
  </Card>
</template>

<style lang="less" scoped></style>
