<script setup lang="ts">
import DomainCountingVue from '@/components/DomainCounting.vue'
import ClientCounting from '@/components/ClientCounting.vue'
import ProcessCounting from './components/ProcessCounting.vue'
import PortCounting from './components/PortCounting.vue'
import ActionBar from './components/ActionBar.vue'
import DNSCounting from './components/DNSCounting.vue'
import DNSTypeCounting from './components/DNSTypeCounting.vue'
import TrafficStatVue from './components/TrafficStat.vue'
import RuleMatchCounting from './components/RuleMatchCounting.vue'
import ProxyCountingVue from './components/ProxyCounting.vue'
</script>

<template>
  <ActionBar></ActionBar>
  <div class="container">
    <div class="left">
      <div class="row1">
        <div class="row1-item">
          <PortCounting />
        </div>
        <div class="row1-item">
          <TrafficStatVue />
        </div>
        <div class="row1-item">
          <ProxyCountingVue />
        </div>
      </div>
      <div class="row2">
        <div class="row2-item">
          <DomainCountingVue />
        </div>
        <div class="row2-item">
          <RuleMatchCounting />
        </div>
      </div>
      <div class="row3">
        <div class="row3-item">
          <ProcessCounting />
        </div>
      </div>
    </div>

    <div class="right">
      <div class="col">
        <div class="col-item">
          <ClientCounting />
        </div>
        <div class="col-item">
          <DNSTypeCounting />
        </div>
        <div class="col-item">
          <DNSCounting />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="less">
.container {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: 8px;
  height: calc(100vh - 46px - 8px);
  .left {
    height: 100%;
    width: 70%;
    display: flex;
    flex-direction: column;
    padding: 0 4px 8px 8px;

    .row1 {
      width: 100%;
      display: flex;
      flex: 1;

      .row1-item {
        margin-left: 8px;
        &:nth-child(1) {
          margin-left: 0;
          width: 25%;
        }
        &:nth-child(2) {
          width: 35%;
        }
        &:nth-child(3) {
          width: 40%;
        }
      }
    }

    .row2 {
      margin-top: 8px;
      width: 100%;
      display: flex;
      flex: 1;

      .row2-item {
        &:nth-child(1) {
          width: 50%;
        }
        &:nth-child(2) {
          margin-left: 8px;
          width: 50%;
        }
      }
    }

    .row3 {
      width: 100%;
      display: flex;
      flex: 1;
      margin-top: 8px;
      .row3-item {
        margin-left: 8px;
        &:nth-child(1) {
          margin-left: 0;
          width: 100%;
        }
        &:nth-child(2) {
          width: 60%;
        }
      }
    }
  }

  .right {
    width: 30%;
    height: 100%;
    padding: 0 8px 8px 4px;
    .col {
      height: 100%;
      display: flex;
      flex-direction: column;
      .col-item {
        margin-top: 8px;
        &:nth-child(1) {
          margin-top: 0;
        }
        &:nth-child(3) {
          flex: 1;
        }
      }
    }
  }
}
</style>
