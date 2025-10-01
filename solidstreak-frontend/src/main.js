import { createApp } from 'vue'
import './style.css'
import App from './App.vue'

createApp(App).mount('#app')

const tg = window.Telegram.WebApp;
console.log('tg.initData:', tg.initData);
console.log('tg.initDataUnsafe:', tg.initDataUnsafe);