<script setup lang='ts'>
import { computed, onMounted, ref } from 'vue'
import { NSpin } from 'naive-ui'
import { fetchChatConfig } from '@/api'
// import pkg from '@/../package.json'
import { useAuthStore } from '@/store'

interface ConfigState {
  timeoutMs?: number
  reverseProxy?: string
  apiModel?: string
  socksProxy?: string
  httpsProxy?: string
  balance?: string
  version?: string
}

const authStore = useAuthStore()

const loading = ref(false)

const config = ref<ConfigState>()

const isChatGPTAPI = computed<boolean>(() => !!authStore.isChatGPTAPI)

async function fetchConfig() {
  try {
    loading.value = true
    const { data } = await fetchChatConfig<ConfigState>()
    config.value = data
  }
  finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchConfig()
})
</script>

<template>
  <NSpin :show="loading">
    <div class="p-4 space-y-4">
      <h2 class="text-xl font-bold">
        <!-- Version - {{ config?.version ?? '-' }} -->
      </h2>
      <div class="p-2 space-y-2 rounded-md bg-neutral-100 dark:bg-neutral-700">
        <p>
          AI影视工业圈是一个聚焦于AI技术与影视制作行业的平台。我们致力于探讨AI技术在影视制作中的实际应用，并关注AI技术对影视产业的革新和改变。同时，我们也会分享关于影视制作中的最新技术和趋势，为从业者和爱好者提供丰富的资源和专业的信息服务。

AI影视工业圈的核心目标是实现AI技术结合影视工业的商业化落地应用。我们鼓励产业内从业者深度参与和研究AI技术，并促进AI技术在影视制作中的不断创新和完善。我们相信，AI技术可以为影视创意孕育更多可能性，并在其中实现更深更广的应用。

AI影视工业圈欢迎具有影视制作背景的专业人士、AI技术研究人员、学者，以及对影视行业未来发展感兴趣的读者加入我们。我们提供行业新闻、深度文章、独家采访和专业咨询等多种资源服务，为您打造更好的影视作品和职业未来。

以上就是AI影视工业圈的星球介绍，希望能够为您提供足够的信息，让更多人加入我们的平台，共同探讨AI技术和影视制作的未来。
        </p>
      </div>
      <p>{{ $t("setting.api") }}：{{ config?.apiModel ?? '-' }}</p>
      <p v-if="isChatGPTAPI">
        {{ $t("setting.balance") }}：{{ config?.balance ?? '-' }}
      </p>
      <p v-if="!isChatGPTAPI">
        {{ $t("setting.reverseProxy") }}：{{ config?.reverseProxy ?? '-' }}
      </p>
      <p>{{ $t("setting.timeout") }}：{{ config?.timeoutMs ?? '-' }}</p>
      <p>{{ $t("setting.socks") }}：{{ config?.socksProxy ?? '-' }}</p>
      <p>{{ $t("setting.httpsProxy") }}：{{ config?.httpsProxy ?? '-' }}</p>
    </div>
  </NSpin>
</template>
