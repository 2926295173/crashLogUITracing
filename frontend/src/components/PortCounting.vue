<script setup lang="ts">
import { computed } from 'vue'
import type { EChartsOption } from 'echarts'
import { useClashTracingStore } from '@/stores'
import BaseCard from './BaseCard.vue'
import BaseChart from './BaseChart.vue'

const store = useClashTracingStore()

const ports = computed(() => store.ruleMatch.portCounting)

const servicesMap: { [key: number]: string } = {
  80: 'HTTP',
  443: 'HTTPS',
  21: 'FTP',
  22: 'SSH',
  23: 'Telnet',
  25: 'SMTP',
  110: 'POP3',
  143: 'IMAP',
  53: 'DNS',
  161: 'SNMP',
  3389: 'RDP',
  5060: 'SIP',
  1433: 'MSSQL',
  3306: 'MySQL'
}

const options = computed(() => {
  const normalPort = ports.value.filter((v) => Object.keys(servicesMap).includes(v.port.toString()))
  const otherPort = ports.value.filter((v) => !Object.keys(servicesMap).includes(v.port.toString()))
  const data = [
    ...normalPort.map(({ port, count }) => ({
      value: count,
      name: servicesMap[port]
    })),
    {
      value: otherPort.reduce((p, c) => p + c.count, 0),
      name: '其他端口：' + otherPort.length + '个'
    }
  ]
  const opt: EChartsOption = {
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'item'
    },
    legend: {
      orient: 'vertical',
      left: 'left'
    },
    series: [
      {
        name: '端口协议',
        type: 'pie',
        radius: '50%',
        data
      }
    ]
  }
  return opt
})
</script>

<template>
  <BaseCard title="端口协议分析">
    <BaseChart :options="options" />
  </BaseCard>
</template>

<style lang="less" scoped></style>
