import { ref } from 'vue'
import { defineStore } from 'pinia'

import * as ApiType from '@/api/types'
import * as ApiFunc from '@/api/index'

const dnsRequest = ref<ApiType.DnsRequestType>({
  counting: [],
  qTypeCounting: [],
  dnsTypeCounting: []
})

const ruleMatch = ref<ApiType.RuleMatchType>({
  counting: [],
  portCounting: [],
  processCounting: [],
  clientCounting: []
})

const proxyDial = ref<ApiType.ProxyDialType>({
  proxyCounting: [],
  hostCounting: []
})

const traffic = ref<ApiType.TrafficType>({
  up: 0,
  down: 0,
  history: []
})

const updateData = async () => {
  const [res1, res2, res3, res4] = await Promise.all([
    ApiFunc.getDnsRequest(),
    ApiFunc.getProxyDial(),
    ApiFunc.getRuleMatch(),
    ApiFunc.getTraffic()
  ])
  dnsRequest.value = res1
  proxyDial.value = res2
  ruleMatch.value = res3
  traffic.value = res4
}

const updateProxyDial = async () => {
  proxyDial.value = await ApiFunc.getProxyDial()
}

const updateRuleMatch = async () => {
  ruleMatch.value = await ApiFunc.getRuleMatch()
}

const updateDnsRequest = async () => {
  dnsRequest.value = await ApiFunc.getDnsRequest()
}

const updateTraffic = async () => {
  traffic.value = await ApiFunc.getTraffic()
}

const syncDB = async () => {
  await ApiFunc.syncDB()
}

export const useClashTracingStore = defineStore('clash-tracing', () => {
  return {
    dnsRequest,
    ruleMatch,
    proxyDial,
    traffic,
    updateData,
    updateProxyDial,
    updateRuleMatch,
    updateDnsRequest,
    updateTraffic,
    syncDB
  }
})
