<script setup lang="ts">
import ClientCounting from './components/ClientCounting.vue'
import DNSCounting from './components/DNSCounting.vue'
import DNSTypeCounting from './components/DNSTypeCounting.vue'
import DomainCountingVue from './components/DomainCounting.vue'
import PortCounting from './components/PortCounting.vue'
import ProcessCounting from './components/ProcessCounting.vue'
import ProxyCountingVue from './components/ProxyCounting.vue'
import RuleMatchCounting from './components/RuleMatchCounting.vue'
import TrafficStatVue from './components/TrafficStat.vue'
</script>

<template>
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
@import url(style.less);
</style>
