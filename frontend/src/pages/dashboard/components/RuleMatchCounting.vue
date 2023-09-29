<script setup lang="ts">
import { computed } from 'vue'
import { useClashTracingStore } from '@/stores'
import Card from '@/components/Card/index.vue'

const store = useClashTracingStore()

const proxies = computed(() => store.ruleMatch.counting)

const maxCountValue = computed(() => Math.max(...proxies.value.map((v) => v.count)))
const maxDurationValue = computed(() => Math.max(...proxies.value.map((v) => v.duration)))

const calcCountWidth = (count: number) => (count / maxCountValue.value) * 100 + '%'
const calcDurationWidth = (duration: number) => (duration / maxDurationValue.value) * 100 + '%'

const formatNumber = (n: number) => n.toFixed(2)
</script>

<template>
  <Card title="规则匹配（次数/速度）">
    <div class="proxies">
      <div v-for="proxy in proxies" :key="proxy.proxy" class="proxies-item hover-item">
        <div class="proxy">{{ proxy.proxy }}({{ proxy.payload || '*' }})</div>
        <div class="proccess-group">
          <div class="processbar">
            <div :style="{ width: calcCountWidth(proxy.count) }" class="process count"></div>
          </div>
          <div class="processbar">
            <div
              :style="{ width: calcDurationWidth(proxy.duration) }"
              class="process duration"
            ></div>
          </div>
        </div>
        <div class="num">
          <div class="count">{{ proxy.count }}</div>
          <div class="duration">{{ formatNumber(proxy.duration / 1000) }} ms</div>
        </div>
      </div>
    </div>
  </Card>
</template>

<style lang="less" scoped>
.proxies {
  &-item {
    display: flex;
    align-items: center;
    font-size: 14px;
    margin: 4px 0;
    .proxy {
      color: #e6e6ed;
      width: 160px;
      text-align: right;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    .proccess-group {
      display: flex;
      flex-direction: column;
      justify-content: center;
      flex: 1;
      border-radius: 4px;
      margin: 0 8px;
      overflow: hidden;
      .processbar {
        height: 10px;
        background: #282b2f;
        .process {
          height: 100%;
          border-radius: 0 4px 4px 0;
        }
        .count {
          background: linear-gradient(to right, #7bc165 70%, #e2575d);
        }
        .duration {
          background: linear-gradient(to right, #642a34 70%, #ec4d5b);
        }
      }
    }
    .num {
      display: flex;
      flex-direction: column;
      text-align: right;
      margin-right: 8px;
      .count {
        color: #73bf69;
        width: 80px;
      }
      .duration {
        color: #be333e;
      }
    }
  }
}
</style>
