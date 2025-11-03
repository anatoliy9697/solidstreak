interface TelegramWebApp {
  initData: string
  onEvent(event: 'webAppReady', callback: () => void): void
}

interface Window {
  Telegram?: {
    WebApp?: TelegramWebApp
  }
}