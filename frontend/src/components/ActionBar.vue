<script setup lang="ts">
import { ref } from 'vue'
import { useClashTracingStore } from '@/stores'
import IconButton from './baseComponents/IconButton.vue'

const title = import.meta.env.VITE_APP_TITLE
const autoRefreshTimer = ref(-1)

const store = useClashTracingStore()

const handleRefresh = () => {
  store.updateData()
}

handleRefresh()

function handleSwitchAutoRefresh() {
  if (autoRefreshTimer.value !== -1) {
    clearInterval(autoRefreshTimer.value)
    autoRefreshTimer.value = -1
    return
  }
  autoRefreshTimer.value = setInterval(() => {
    store.updateData()
  }, 5000)
}

async function handleSyncDB() {
  await store.syncDB()
  store.updateData()
}
</script>

<template>
  <div class="actionbar">
    <div class="logo">{{ title }}</div>
    <div class="action">
      <IconButton @click="handleRefresh" icon="refresh" title="刷新"> </IconButton>
      <IconButton
        @click="handleSwitchAutoRefresh"
        text="自动刷新"
        :icon="autoRefreshTimer === -1 ? 'close' : 'open'"
      >
      </IconButton>
      <IconButton @click="handleSyncDB" title="同步内存的数据到数据库" icon="sync" />
      <IconButton
        icon="github"
        type="link"
        href="https://github.com/openrhc/Clash-Tracing"
        style="margin-left: 4px"
      />
      <!-- <IconButton text="Last 1 hour" icon="dropdown" style="margin-left: 4px">
        <div class="select-time">
          <div value="1">1</div>
          <div value="2">2</div>
          <div value="3">3</div>
        </div>
      </IconButton> -->
    </div>
  </div>
</template>

<style lang="less" scoped>
.actionbar {
  display: flex;
  align-items: center;
  height: 46px;
  padding: 0 16px;
  background: #181b1e;
  .logo {
    color: #fff;
    font-size: 18px;
    font-weight: bold;
  }
  .action {
    margin-left: auto;
    display: flex;
    align-items: center;

    .select-time {
      color: #fff;
      div {
        padding: 0 8px;
        line-height: 36px;
        &:hover {
          background: #181b1e;
        }
      }
    }
  }
}
</style>
