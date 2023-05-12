import { ss } from '@/utils/storage'

const LOCAL_NAME = 'userStorage'

export interface UserInfo {
  avatar: string
  name: string
  description: string
  baseURI: string
  accessToken: string
}

export interface UserState {
  userInfo: UserInfo
}

export function defaultSetting(): UserState {
  // 生成随机的字符串
  function generateRandomString(length: number): string {
    let result = ''
    const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
    for (let i = 0; i < length; i++)
      result += characters.charAt(Math.floor(Math.random() * characters.length))

    return result
  }
  const randomString = generateRandomString(10)
  return {
    userInfo: {
      avatar: 'https://xyhelper.cn/defaultavatar.jpeg',
      name: '攻城狮老李',
      description: '访问 <a href="https://xyhelper.cn" class="text-blue-500" target="_blank" >xyhelper.cn</a>',
      baseURI: 'https://freechat.xyhelper.cn',
      accessToken: randomString,
    },
  }
}

export function getLocalState(): UserState {
  const localSetting: UserState | undefined = ss.get(LOCAL_NAME)
  return { ...defaultSetting(), ...localSetting }
}

export function setLocalState(setting: UserState): void {
  ss.set(LOCAL_NAME, setting)
}
