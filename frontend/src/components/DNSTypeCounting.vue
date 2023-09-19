<script setup lang="ts">
import { computed } from 'vue'
import type { EChartsOption } from 'echarts'
import { useClashTracingStore } from '@/stores'
import BaseCard from './BaseCard.vue'
import BaseChart from './BaseChart.vue'

const store = useClashTracingStore()

const dnsType = computed(() => store.dnsRequest.dnsTypeCounting)

const dnsQType = computed(() => store.dnsRequest.qTypeCounting)

const dnsTypeOption = computed(() => {
  const opt: EChartsOption = {
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'item'
    },
    series: [
      {
        name: '端口协议',
        type: 'pie',
        radius: '50%',
        data: dnsType.value.map(({ dnsType, count }) => ({
          value: count,
          name: dnsType
        }))
      }
    ]
  }
  return opt
})

const dnsQTypeOption = computed(() => {
  const opt: EChartsOption = {
    backgroundColor: 'transparent',
    xAxis: {
      type: 'category',
      data: dnsQType.value.map((v) => v.qType)
    },
    grid: {
      left: '0%',
      right: '0%',
      bottom: '0%',
      top: '8%',
      containLabel: true
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        data: dnsQType.value.map((v) => v.count),
        type: 'bar',
        label: {
          formatter: (e) => {
            console.log(e)

            return '123'
          }
        },
        showBackground: true,
        backgroundStyle: {
          color: 'rgba(180, 180, 180, 0.2)'
        }
      }
    ]
  }
  return opt
})
</script>

<template>
  <BaseCard title="DNS查询类型" min-height="170px">
    <div style="display: flex">
      <BaseChart height="160px" :options="dnsTypeOption" />
      <BaseChart height="160px" :options="dnsQTypeOption" />
    </div>
  </BaseCard>
</template>

<style lang="less" scoped></style>
