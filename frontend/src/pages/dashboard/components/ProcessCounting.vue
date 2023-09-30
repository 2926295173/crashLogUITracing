<script setup lang="ts">
import { computed, ref } from 'vue'
import { useClashTracingStore } from '@/stores'
import { useState } from '@/hooks/useState'
import { getProcessDetail } from '@/api/index'
import { formatTime } from '@/utils'
import Card from '@/components/Card/index.vue'
import Modal from '@/components/Modal/index.vue'
import Table, { type TableColumnItem } from '@/components/Table/index.vue'

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
    render: ({ val }) => formatTime(val)
  }
]

const processPath = ref('')
const processes = computed(() => store.ruleMatch.processCounting)

const store = useClashTracingStore()
const [isModalOpen, setIsModalOpen] = useState(false)

const handleDetail = async (path: string) => {
  processPath.value = path
  setIsModalOpen(true)
}
</script>

<template>
  <Card title="进程统计">
    <div class="processes">
      <div
        v-for="process in processes"
        :key="process.path"
        @click="handleDetail(process.path)"
        class="processes-item hover-item"
      >
        <div class="path">{{ process.path || 'localhost' }}</div>
        <div class="count">{{ process.count }}</div>
      </div>
    </div>
  </Card>
  <Modal :open="isModalOpen" @on-cancel="setIsModalOpen(false)">
    <Card title="进程详情" min-height="450px">
      <Table
        :columns="columns"
        :params="{ path: processPath }"
        :title="processPath"
        :api="getProcessDetail"
      />
    </Card>
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
