<script setup lang="ts">
import { computed } from 'vue'
import { useClashTracingStore } from '@/stores'
import BaseCard from './BaseCard.vue'

const store = useClashTracingStore()

const processes = computed(() =>
  store.ruleMatch.processCounting.map((v) => ({ ...v, path: v.path || 'localhost' }))
)
</script>

<template>
  <BaseCard title="进程统计">
    <div class="processes">
      <div v-for="process in processes" :key="process.path" class="processes-item hover-item">
        <div class="path">{{ process.path }}</div>
        <div class="count">{{ process.count }}</div>
      </div>
    </div>
  </BaseCard>
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
