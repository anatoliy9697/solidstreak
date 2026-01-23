interface WebAppUser {
  id: number
  username?: string
  first_name: string
  last_name?: string
  language_code?: string
  is_bot?: boolean
  photo_url?: string
}

interface WebAppChat {
  id: number
}

interface WebAppInitData {
  user?: WebAppUser
  chat?: WebAppChat
}

interface TelegramWebApp {
  initData: string
  initDataUnsafe: WebAppInitData
  onEvent(event: 'webAppReady', callback: () => void): void
  ready(): void
}

interface Window {
  Telegram?: {
    WebApp?: TelegramWebApp
  }
}
