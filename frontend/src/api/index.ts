import type { AxiosProgressEvent, GenericAbortSignal } from 'axios'
import { post } from '@/utils/request'
import { useSettingStore } from '@/store'

export function fetchChatAPI<T = any>(
  prompt: string,
  options?: { conversationId?: string; parentMessageId?: string },
  signal?: GenericAbortSignal,
) {
  return post<T>({
    url: '/chat',
    data: { prompt, options },
    signal,
  })
}

export function fetchChatConfig<T = any>() {
  return post<T>({
    url: '/config',
  })
}

export function fetchChatAPIProcess<T = any>(
  params: {
    prompt: string
    options?: { conversationId?: string; parentMessageId?: string }
    baseURI?: string
    accessToken?: string
    isGPT4?: boolean
    signal?: GenericAbortSignal
    onDownloadProgress?: (progressEvent: AxiosProgressEvent) => void },
) {
  const settingStore = useSettingStore()

  return post<T>({
    url: '/chat-process',
    data: { prompt: params.prompt, options: params.options, systemMessage: settingStore.systemMessage, baseURI: params.baseURI, accessToken: params.accessToken, isGPT4: params.isGPT4 },
    signal: params.signal,
    onDownloadProgress: params.onDownloadProgress,
  })
}

export function fetchSession<T>() {
  return post<T>({
    url: '/session',
  })
}

export function fetchVerify<T>(token: string) {
  return post<T>({
    url: '/verify',
    data: { token },
  })
}
