<script setup lang="ts">
import { computed, ref } from 'vue'
import { useClashTracingStore } from '@/stores'
import { useState } from '@/hooks/useState'
import { getProcessDetail } from '@/api/index'
import type { RuleMatchType, ProcessDetailType } from '@/api/types'
import { formatTime } from '@/utils'
import BaseCard from './BaseCard.vue'
import Modal from './baseComponents/Modal.vue'
import Table, { type TableColumnItem } from './baseComponents/Table.vue'

const dataSource = ref<ProcessDetailType>([])
const columns: TableColumnItem[] = [
  {
    title: '源IP',
    dataIndex: 'sourceIP'
  },
  {
    title: '源Port',
    dataIndex: 'sourcePort'
  },
  {
    title: '目标IP',
    dataIndex: 'destinationIP'
  },
  {
    title: '目标Port',
    dataIndex: 'destinationPort'
  },
  {
    title: '主机',
    dataIndex: 'host'
  },
  {
    title: 'dns模式',
    dataIndex: 'dnsMode'
  },
  {
    title: '时间',
    dataIndex: 'createTime',
    render: (val: number) => formatTime(val)
  }
]

const store = useClashTracingStore()
const [isModalOpen, setIsModalOpen] = useState(false)

const processes = computed(() =>
  store.ruleMatch.processCounting.map((v) => ({ ...v, path: v.path || 'localhost' }))
)

const handleDetail = async ({ path }: RuleMatchType['processCounting'][0]) => {
  dataSource.value = []
  setIsModalOpen(true)
  dataSource.value = await getProcessDetail(path)
}
</script>

<template>
  <BaseCard title="进程统计">
    <div class="processes">
      <div
        v-for="process in processes"
        :key="process.path"
        @click="handleDetail(process)"
        class="processes-item hover-item"
      >
        <div class="path">{{ process.path }}</div>
        <div class="count">{{ process.count }}</div>
      </div>
    </div>
  </BaseCard>
  <Modal :open="isModalOpen" @on-cancel="setIsModalOpen(false)">
    <BaseCard title="进程详情" min-height="220px">
      <Table :columns="columns" :data-source="dataSource" />
    </BaseCard>
  </Modal>
</template>

<style lang="less" scoped>
.processes {
  &-item {
    display: flex;
    align-items: center;
    font-size: 14px;
    margin: 4px 0;

    .path {
      color: #e6e6ed;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    .count {
      color: #73bf69;
      text-align: right;
      margin-left: auto;
      margin-right: 8px;
      width: 80px;
    }
  }
}
</style>
