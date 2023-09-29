<script setup lang="ts">
import { computed } from 'vue'
import { useClashTracingStore } from '@/stores'
import Card from '@/components/Card/index.vue'

const store = useClashTracingStore()

const domains = computed(() => store.proxyDial.hostCounting)

const maxValue = computed(() => Math.max(...domains.value.map((v) => v.count)))

const calcWidth = (count: number) => (count / maxValue.value) * 100 + '%'
</script>

<template>
  <Card title="域名访问次数">
    <div class="domains">
      <div v-for="domain in domains" :key="domain.host" class="domains-item hover-item">
        <div class="domain">{{ domain.host }}</div>
        <div class="processbar">
          <div :style="{ width: calcWidth(domain.count) }" class="process"></div>
        </div>
        <div class="count">{{ domain.count }}</div>
      </div>
    </div>
  </Card>
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
