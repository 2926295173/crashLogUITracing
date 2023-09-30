<script setup lang="ts" name="Table">
import { computed, ref, watch } from 'vue'

type DataSourceItem = Record<string, any>

export type TableColumnItem = {
  title: string
  dataIndex: string
  render?: (val: { record: DataSourceItem; val: any }) => any
  width?: number
}

interface Props {
  title?: string
  columns: TableColumnItem[]
  params?: Record<string, any>
  api?: (params: Record<string, any>) => Promise<{
    total: number
    page: number
    pageSize: number
    data: DataSourceItem[]
  }>
}

const pagerSize = 6

const page = ref<number>(1)
const pageSize = ref<number>(10)
const total = ref<number>(0)
const dataSource = ref<DataSourceItem[]>([])
const pageTotal = computed(() => Math.ceil(total.value / pageSize.value))

const preList = computed(() =>
  new Array(pagerSize)
    .fill(page.value)
    .map((i, idx) => i - idx - 1)
    .filter((v) => v > 0)
    .reverse()
)

const nextList = computed(() =>
  new Array(pagerSize)
    .fill(page.value)
    .map((i, idx) => i + idx - 1)
    .filter((v) => v > page.value)
    .filter((v) => v <= pageTotal.value)
)

const props = withDefaults(defineProps<Props>(), {
  title: '',
  params: () => ({}),
  api: async () => ({
    total: 0,
    page: 1,
    pageSize: 10,
    data: []
  })
})

const fetchData = async () => {
  const d = await props.api({
    ...props.params,
    page: page.value,
    pageSize: pageSize.value
  })
  total.value = d.total
  page.value = d.page
  pageSize.value = d.pageSize
  dataSource.value = d.data
}

const changePage = (p: number) => {
  page.value = p
  fetchData()
}

watch(() => props.api, fetchData, { immediate: true })
</script>

<template>
  <div v-if="props.title" class="title">{{ props.title }}</div>
  <table border="0">
    <thead>
      <tr>
        <th
          v-for="c in columns"
          :key="c.dataIndex"
          :style="{ width: c.width ? c.width + 'px' : 'auto' }"
        >
          {{ c.title }}
        </th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="(d, idx) in dataSource" :key="idx" class="hover-item">
        <td
          v-for="c in columns"
          :key="c.dataIndex + idx"
          :style="{ width: c.width ? c.width + 'px' : 'auto' }"
        >
          {{ c.render ? c.render({ record: d, val: d[c.dataIndex] }) : d[c.dataIndex] }}
        </td>
      </tr>
    </tbody>
  </table>
  <div v-show="pageTotal > 1" class="pager">
    <div v-show="preList.length !== 0" @click="changePage(page - 1)" class="pager-item">上一页</div>
    <div v-for="p in preList" :key="p" @click="changePage(p)" class="pager-item">{{ p }}</div>
    <div class="pager-item page">{{ page }}</div>
    <div v-for="p in nextList" :key="p" @click="changePage(p)" class="pager-item">{{ p }}</div>
    <div v-show="nextList.length !== 0" @click="changePage(page + 1)" class="pager-item">
      下一页
    </div>
  </div>
</template>

<style lang="less" scoped>
@import url(style.less);
</style>
