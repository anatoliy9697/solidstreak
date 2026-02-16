import './style.css'
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import PrimeVue from 'primevue/config'
import Aura from '@primeuix/themes/aura'
import ToastService from 'primevue/toastservice'
import ConfirmationService from 'primevue/confirmationservice'
import { i18n } from '@/i18n'
import { type Theme, THEMES, toLocalTheme } from '@/models/theme'

import App from './App.vue'
// import router from './router'

export function applyTheme(theme: Theme) {
  const root = document.documentElement
  const themeVars = THEMES[theme]
  Object.entries(themeVars).forEach(([key, value]) => {
    root.style.setProperty(key, value)
  })
  if (theme === 'dark') {
    root.classList.add('dark')
  } else {
    root.classList.remove('dark')
  }
}

applyTheme(toLocalTheme(localStorage.getItem('theme')))

const app = createApp(App)

app.use(i18n)
// app.use(router)
app.use(PrimeVue, {
  theme: {
    preset: Aura,
    options: {
      darkModeSelector: false || 'none',
    },
  },
})
app.use(createPinia())
app.use(ToastService)
app.use(ConfirmationService)

app.mount('#app')
