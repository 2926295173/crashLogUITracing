<script setup lang="ts">
import { computed, ref } from 'vue'
import { useClashTracingStore } from '@/stores'
import { useState } from '@/hooks/useState'
import { getDomainDetail } from '@/api/index'
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
    title: '网络',
    dataIndex: 'network'
  },
  {
    title: 'dns模式',
    dataIndex: 'dnsMode'
  },
  {
    title: '入口',
    dataIndex: 'type'
  },
  {
    title: '规则代理',
    dataIndex: 'rule',
    render: ({ record, val }) => {
      return val + '::' + record.payload + '->' + record.proxy
    }
  },
  {
    title: '进程',
    dataIndex: 'processPath'
  },
  {
    title: '时间',
    dataIndex: 'createTime',
    render: ({ val }) => formatTime(val)
  }
]

const host = ref('')
const store = useClashTracingStore()
const [isModalOpen, setIsModalOpen] = useState(false)

const domains = computed(() => store.proxyDial.hostCounting)

const maxValue = computed(() => Math.max(...domains.value.map((v) => v.count)))

const calcWidth = (count: number) => (count / maxValue.value) * 100 + '%'

const handleDetail = async (_host: string) => {
  host.value = _host
  setIsModalOpen(true)
}
</script>

<template>
  <Card title="域名访问次数">
    <div class="domains">
      <div
        v-for="domain in domains"
        :key="domain.host"
        @click="handleDetail(domain.host)"
        class="domains-item hover-item"
      >
        <div class="domain">{{ domain.host }}</div>
        <div class="processbar">
          <div :style="{ width: calcWidth(domain.count) }" class="process"></div>
        </div>
        <div class="count">{{ domain.count }}</div>
      </div>
    </div>
  </Card>
  <Modal :open="isModalOpen" @on-cancel="setIsModalOpen(false)">
    <Card title="请求详情" min-height="450px">
      <Table :columns="columns" :params="{ host: host }" :title="host" :api="getDomainDetail" />
    </Card>
  </Modal>
</template>

<style lang="less" scoped>
.domains {
  &-item {
    display: flex;
    align-items: center;
    font-size: 14px;
    margin: 4px 0;
    .domain {
      color: #e6e6ed;
      width: 220px;
      text-align: right;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    .processbar {
      flex: 1;
      margin: 0 8px;
      height: 10px;
      background: #282b2f;
      border-radius: 4px;
      .process {
        height: 100%;
        border-radius: 4px;
        background: linear-gradient(to right, #73bf69 70%, #ec4d5b);
      }
    }
    .count {
      color: #73bf69;
      margin-right: 8px;
      width: 80px;
      text-align: right;
    }
  }
}
</style>
