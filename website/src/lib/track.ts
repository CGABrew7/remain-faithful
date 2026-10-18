const sent = new Set<string>()

export function trackingOff(): boolean {
  if (typeof navigator === 'undefined') return true
  const nav = navigator as Navigator & { msDoNotTrack?: string; globalPrivacyControl?: boolean }
  const win = window as Window & { doNotTrack?: string }
  const dnt = nav.doNotTrack || win.doNotTrack || nav.msDoNotTrack
  if (dnt === '1' || dnt === 'yes') return true
  if (nav.globalPrivacyControl === true) return true
  return false
}

export function trackEvent(
  event: 'pageview' | 'app_store_click' | 'android_signup',
  path?: string,
) {
  if (typeof window === 'undefined' || trackingOff()) return
  const page = path || window.location.pathname
  if (page === '/admin' || page.startsWith('/admin/')) return
  if (event === 'pageview') {
    const key = `${page}?${window.location.search}`
    if (sent.has(key)) return
    sent.add(key)
  }
  const params = new URLSearchParams(window.location.search)
  const body = JSON.stringify({
    event,
    path: page,
    referrer: document.referrer || '',
    utm_source: params.get('utm_source') || '',
    utm_medium: params.get('utm_medium') || '',
    utm_campaign: params.get('utm_campaign') || '',
  })
  const blob = new Blob([body], { type: 'application/json' })
  if (navigator.sendBeacon?.('/api/collect', blob)) return
  void fetch('/api/collect', {
    method: 'POST',
    body,
    headers: { 'Content-Type': 'application/json' },
    keepalive: true,
  })
}
