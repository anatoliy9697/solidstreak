import { defineStore } from 'pinia'

import { ApiFetcher, type RequestResult } from '@/api/request'
import type { User } from '@/models/user'

export const useUserStore = defineStore('user', {
  state: () => ({
    apiFetcher: null as ApiFetcher | null,
    id: 0 as number,
    tgId: 0 as number,
    tgUsername: '' as string,
    tgFirstName: '' as string,
    tgLastName: '' as string,
    tgLangCode: '' as string,
    avatarUrl: '' as string,
  }),

  actions: {
    init(apiFetcher: ApiFetcher): void {
      this.apiFetcher = apiFetcher
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
      }

      return result
    },
  },
})
