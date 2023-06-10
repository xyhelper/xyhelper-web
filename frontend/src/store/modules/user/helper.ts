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
  // 获取 https://freechat.xyhelper.cn/default-setting.json 中的配置
  // async function getRemoteSetting(): Promise<UserState> {
  //   const response = await fetch('https://freechat.xyhelper.cn/default-setting.json')
  //   const data = await response.json()
  //   return data
  // }
  // 从 https://freechat.xyhelper.cn/default-setting.json 中获取配置
  // 如果获取失败，则使用默认配置
  // try {
  //   const remoteSetting = getRemoteSetting()
  //   return remoteSetting as unknown as UserState
  // }
  // catch (error) {
  // console.log(error)
  return {
    userInfo: {
      avatar: '',
      name: '人类',
      description: '分享AI知识，共享AI资源',
      baseURI: 'https://pluschat.xyhelper.com.cn',
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
