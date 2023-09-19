<script setup lang="ts">
import { computed } from 'vue'
import { useClashTracingStore } from '@/stores'
import BaseCard from './BaseCard.vue'

const store = useClashTracingStore()

const dns = computed(() => store.dnsRequest.counting)

const maxValue = computed(() => Math.max(...dns.value.map((v) => v.count)))

const calcWidth = (count: number) => (count / maxValue.value) * 100 + '%'
</script>

<template>
  <BaseCard title="DNS查询次数">
    <div class="dns">
      <div v-for="d in dns" :key="d.name" class="dns-item hover-item">
        <div class="domain">{{ d.name }}</div>
        <div class="processbar">
          <div :style="{ width: calcWidth(d.count) }" class="process"></div>
        </div>
        <div class="count">{{ d.count }}</div>
      </div>
    </div>
  </BaseCard>
</template>

<style lang="less" scoped>
.dns {
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
