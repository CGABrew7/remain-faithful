'use client'

import { useCallback, useEffect, useId, useRef } from 'react'
import Script from 'next/script'

type TurnstileWidgetId = string

type TurnstileApi = {
  render: (
    container: HTMLElement,
    options: {
      sitekey: string
      action: string
      appearance: 'interaction-only'
      execution: 'render'
      theme: 'auto'
      callback: (token: string) => void
      'expired-callback': () => void
      'error-callback': () => void
    },
  ) => TurnstileWidgetId
  reset: (widgetId: TurnstileWidgetId) => void
  remove: (widgetId: TurnstileWidgetId) => void
}

declare global {
  interface Window {
    turnstile?: TurnstileApi
  }
}

export type FormShieldFields = {
  company: string
  turnstileToken: string
  loadedAt: number
  formAction: string
}

export function formShieldSiteKey() {
  return process.env.NEXT_PUBLIC_TURNSTILE_SITE_KEY || ''
}

// useFormShield adds a hidden honeypot and a managed Turnstile widget.
// The widget stays invisible unless Cloudflare asks for an interaction.
export function useFormShield(action: string) {
  const siteKey = formShieldSiteKey()
  const loadedAt = useRef(0)
  if (typeof window !== 'undefined' && loadedAt.current === 0) {
    loadedAt.current = Date.now()
  }
  const tokenRef = useRef('')
  const companyRef = useRef<HTMLInputElement>(null)
  const widgetId = useRef<TurnstileWidgetId | null>(null)
  const pollRef = useRef<number | null>(null)
  const fieldId = useId()

  const clearPoll = useCallback(() => {
    if (pollRef.current !== null) {
      window.clearInterval(pollRef.current)
      pollRef.current = null
    }
  }, [])

  const removeWidget = useCallback(() => {
    clearPoll()
    const id = widgetId.current
    widgetId.current = null
    tokenRef.current = ''
    if (id && window.turnstile) window.turnstile.remove(id)
  }, [clearPoll])

  useEffect(() => {
    return () => {
      removeWidget()
    }
  }, [removeWidget])

  const containerRef = useCallback(
    (node: HTMLDivElement | null) => {
      if (!node) {
        removeWidget()
        return
      }
      clearPoll()
      if (!siteKey) return
      const render = () => {
        if (widgetId.current !== null || !window.turnstile) return false
        widgetId.current = window.turnstile.render(node, {
          sitekey: siteKey,
          action,
          appearance: 'interaction-only',
          execution: 'render',
          theme: 'auto',
          callback: (token: string) => {
            tokenRef.current = token
          },
          'expired-callback': () => {
            tokenRef.current = ''
          },
          'error-callback': () => {
            tokenRef.current = ''
          },
        })
        return true
      }
      if (render()) return
      pollRef.current = window.setInterval(() => {
        if (render()) clearPoll()
      }, 200)
    },
    [action, clearPoll, removeWidget, siteKey],
  )

  const fields = useCallback((): FormShieldFields => {
    return {
      company: companyRef.current?.value || '',
      turnstileToken: tokenRef.current,
      loadedAt: loadedAt.current,
      formAction: action,
    }
  }, [action])

  const reset = useCallback(() => {
    tokenRef.current = ''
    if (widgetId.current !== null && window.turnstile) {
      window.turnstile.reset(widgetId.current)
    }
  }, [])

  const slot = (
    <div className="pointer-events-none absolute left-0 top-0 z-10 h-0 w-0 overflow-visible">
      <div className="absolute left-[-9999px] h-0 w-0 overflow-hidden" aria-hidden="true">
        <label htmlFor={fieldId}>Company</label>
        <input
          ref={companyRef}
          id={fieldId}
          type="text"
          name="company"
          tabIndex={-1}
          autoComplete="off"
          defaultValue=""
        />
      </div>
      <div ref={containerRef} className="pointer-events-auto absolute left-0 top-0" />
      {siteKey ? (
        <Script
          src="https://challenges.cloudflare.com/turnstile/v0/api.js?render=explicit"
          strategy="afterInteractive"
        />
      ) : null}
    </div>
  )

  return { fields, reset, slot }
}
