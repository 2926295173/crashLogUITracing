const apiUrl = import.meta.env.VITE_API_URL || ''

export const httpGet = async <T>(url: string, params: Record<string, any> = {}) => {
  const usp = new URLSearchParams()
  Object.entries(params).forEach(([key, val]) => {
    usp.append(key, val)
  })
  let _url = apiUrl + url
  if (usp.toString()) {
    _url += '?' + usp.toString()
  }
  const res = await fetch(_url)
  const result = await res.json()
  return result as T
}
