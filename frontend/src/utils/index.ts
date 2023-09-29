export const formatTime = (smillisecond: number) => {
  const d = new Date(smillisecond * 1000)
  const h = d.getHours().toString().padStart(2, '0')
  const m = d.getMinutes().toString().padStart(2, '0')
  const s = d.getSeconds().toString().padStart(2, '0')
  return `${h}:${m}:${s}`
}

export function formatTraffic(num: number) {
  const s = ['B', 'KiB', 'MiB', 'GiB', 'TiB']
  const exp = Math.floor(Math.log(num || 1) / Math.log(1024))
  return `${(Math.floor((num / Math.pow(1024, exp)) * 100) / 100).toFixed(2)} ${s?.[exp] ?? ''}`
}
