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
      avatar: 'avatar.jpg',
      name: '米卡传媒',
      description: 'AI技术与影视制作平台',
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
