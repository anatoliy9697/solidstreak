import { access, mkdir, writeFile } from 'node:fs/promises'
import { constants } from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const SDK_URL = 'https://telegram.org/js/telegram-web-app.js'
const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)
const projectRoot = path.resolve(__dirname, '..')
const targetDir = path.join(projectRoot, 'public', 'vendor')
const targetFile = path.join(targetDir, 'telegram-web-app.js')

async function exists(filePath) {
  try {
    await access(filePath, constants.F_OK)
    return true
  } catch {
    return false
  }
}

async function main() {
  await mkdir(targetDir, { recursive: true })

  try {
    const response = await fetch(SDK_URL, {
      headers: { 'User-Agent': 'solidstreak-sdk-sync' },
    })

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`)
    }

    const body = await response.text()

    await writeFile(targetFile, body, 'utf8')
    console.log(`Telegram SDK synchronized: ${targetFile}`)
  } catch (err) {
    console.error(
      `Failed to synchronize Telegram SDK: ${err instanceof Error ? err.message : String(err)}`,
    )
    process.exit(1)
  }
}

await main()