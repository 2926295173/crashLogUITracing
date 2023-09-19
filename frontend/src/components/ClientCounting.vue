<script setup lang="ts">
import { computed } from 'vue'
import { useClashTracingStore } from '@/stores'
import BaseCard from './BaseCard.vue'

const store = useClashTracingStore()

const clients = computed(() => store.ruleMatch.clientCounting)

const maxValue = computed(() => Math.max(...clients.value.map((v) => v.count)))

const calcWidth = (count: number) => (count / maxValue.value) * 100 + '%'
</script>

<template>
  <BaseCard title="局域网访问" min-height="130px">
    <div class="clients">
      <div v-for="client in clients" :key="client.ip" class="clients-item hover-item">
        <div class="ip">{{ client.ip }}</div>
        <div class="processbar">
          <div :style="{ width: calcWidth(client.count) }" class="process"></div>
        </div>
        <div class="count">{{ client.count }}</div>
      </div>
    </div>
  </BaseCard>
</template>

<style lang="less" scoped>
.clients {
  &-item {
    display: flex;
    align-items: center;
    font-size: 14px;
    margin: 4px 0;

    .ip {
      color: #73bf69;
      width: 120px;
      font-size: 16px;
      text-align: right;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    .processbar {
      flex: 1;
      margin: 4px 8px;
      height: 26px;
      background: #282b2f;
      border-radius: 4px;
      .process {
        height: 100%;
        border-radius: 4px;
        background: linear-gradient(to right, #7ac165, #face2f, #f3515a);
      }
    }
    .count {
      color: #73bf69;
      font-size: 22px;
      text-align: right;
      margin-right: 8px;
      width: 80px;
      &:nth-child(1) {
        color: #e44658;
      }
      &:nth-child(2) {
        color: #dfc33e;
      }
    }
  }
}
</style>
