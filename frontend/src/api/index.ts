import { httpGet } from '@/utils/request'
import type {
  DnsRequestType,
  RuleMatchType,
  ProxyDialType,
  TrafficType,
  ProcessDetailType,
  PageType
} from './types'

enum Api {
  DnsRequest = '/dnsRequest',
  RuleMatch = '/ruleMatch',
  ProxyDial = '/proxyDial',
  Traffic = '/traffic',
  Sync = '/sync',
  ProcessDetail = '/processDetail'
}

export const getDnsRequest = () => httpGet<DnsRequestType>(Api.DnsRequest)

export const getRuleMatch = () => httpGet<RuleMatchType>(Api.RuleMatch)

export const getProxyDial = () => httpGet<ProxyDialType>(Api.ProxyDial)

export const getTraffic = () => httpGet<TrafficType>(Api.Traffic)

export const syncDB = () => httpGet<null>(Api.Sync)

export const getProcessDetail = (params: Record<string, any>) =>
  httpGet<PageType<ProcessDetailType>>(Api.ProcessDetail, params)
