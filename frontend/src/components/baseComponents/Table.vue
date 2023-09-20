<script setup lang="ts">
export type TableColumnItem = {
  title: string
  dataIndex: string
  render?: (val: any) => any
}

interface Props {
  dataSource: Record<string, any>[]
  columns: TableColumnItem[]
}

defineProps<Props>()
</script>

<template>
  <table>
    <thead>
      <tr>
        <th v-for="c in columns" :key="c.dataIndex">{{ c.title }}</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="(d, idx) in dataSource" :key="idx">
        <td v-for="c in columns" :key="c.dataIndex + idx">
          {{ c.render ? c.render(d[c.dataIndex]) : d[c.dataIndex] }}
        </td>
      </tr>
    </tbody>
  </table>
</template>

<style lang="less" scoped>
table {
  text-align: center;
  font-size: 14px;
  tbody {
    overflow-y: auto;
  }
  tr {
    line-height: 24px;
  }
}
</style>
