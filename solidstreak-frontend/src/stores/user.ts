import { defineStore } from 'pinia'

import { toLocalLang } from '@/i18n'
import { ApiFetcher, type RequestResult } from '@/api/request'
import type { User } from '@/models/user'
import type { Subscription } from '@/models/subscription'
import { type Theme, getDefaultTheme, toLocalTheme } from '@/models/theme'

export const useUserStore = defineStore('user', {
  state: () => ({
    apiFetcher: null as ApiFetcher | null,
    id: 0 as number,
    tgId: 0 as number,
    tgUsername: '' as string,
    tgFirstName: '' as string,
    tgLastName: '' as string,
    tgLangCode: '' as string,
    langCode: toLocalLang(localStorage.getItem('lang')) as string,
    avatarUrl: '' as string,
    _subscription: undefined as Subscription | undefined,
    _theme: toLocalTheme(localStorage.getItem('theme') || getDefaultTheme()) as Theme,
  }),

  actions: {
    init(apiFetcher: ApiFetcher, defaultSubscription: Subscription): void {
      this.apiFetcher = apiFetcher
      this._subscription = defaultSubscription
    },

    setAvatarUrl(avatarUrl: string): void {
      this.avatarUrl = avatarUrl
    },

    async upsertUserInfo(webAppUser: WebAppUser, webAppChat: WebAppChat): Promise<RequestResult> {
      const inputUser = {
        tgId: webAppUser.id,
        tgUsername: webAppUser.username,
        tgFirstName: webAppUser.first_name,
        tgLastName: webAppUser.last_name,
        tgLangCode: webAppUser.language_code,
        langCode: toLocalLang(localStorage.getItem('lang') || webAppUser.language_code),
        tgIsBot: webAppUser.is_bot,
      } as User

      const result = await this.apiFetcher!.upsertUserInfo(inputUser, { tgId: webAppChat.id })

      const user = result.response?.data ? (result.response?.data as User) : null

      if (user) {
        this.id = user.id || 0
        this.tgId = user.tgId
        this.tgUsername = user.tgUsername || ''
        this.tgFirstName = user.tgFirstName
        this.tgLastName = user.tgLastName || ''
        this.tgLangCode = user.tgLangCode || ''
        this.langCode = toLocalLang(
          user.langCode || localStorage.getItem('lang') || user.tgLangCode,
        )
        localStorage.setItem('lang', this.langCode)
        this._subscription = user.subscription || this._subscription
        if (this._subscription!.finishDate)
          this._subscription!.finishDate = new Date(this._subscription!.finishDate) // Преобразуем в Date, т.к. по факту из API оно приходит в строкой
      }

      return result
    },

    async setLang(lang: string): Promise<RequestResult> {
      const prevLang = this.langCode

      this.langCode = lang
      localStorage.setItem('lang', lang)

      const result = await this.apiFetcher!.patchUser(this.id, lang)

      if (!result.success) {
        this.langCode = prevLang
        localStorage.setItem('lang', prevLang)
      }

      return result
    },

    setTheme(theme: Theme): void {
      this._theme = theme
      localStorage.setItem('theme', theme)
    },
  },

  getters: {
    lang: (state): string => {
      return toLocalLang(state.langCode || localStorage.getItem('lang') || state.tgLangCode)
    },

    theme: (state): Theme => {
      if (!localStorage.getItem('theme')) {
        localStorage.setItem('theme', state._theme || getDefaultTheme())
      }
      return toLocalTheme(state._theme || localStorage.getItem('theme') || getDefaultTheme())
    },

    subscription: (state): Subscription => {
      return state._subscription!
    },

    activeHabitsLimit: (state): number => {
      return state._subscription!.plan.activeHabitsLimit
    },
  },
})
