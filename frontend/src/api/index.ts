import { httpGet } from '@/utils/request'
import type {
  DnsRequestType,
  RuleMatchType,
  ProxyDialType,
  TrafficType,
  ProcessDetailType
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

export const getProcessDetail = (path: string) =>
  httpGet<ProcessDetailType>(Api.ProcessDetail, { path })
