<script lang="ts" setup>
import { computed } from 'vue'
import { NAvatar } from 'naive-ui'
import { useUserStore } from '@/store'
import { isString } from '@/utils/is'
import defaultAvatar from '@/assets/avatar.jpg'

interface Props {
  image?: boolean
}
defineProps<Props>()

const userStore = useUserStore()

const avatar = computed(() => userStore.userInfo.avatar)
</script>

<template>
  <template v-if="image">
    <NAvatar v-if="isString(avatar) && avatar.length > 0" :src="avatar" :fallback-src="defaultAvatar" />
    <NAvatar v-else round :src="defaultAvatar" />
  </template>
  <span v-else class="text-[28px] dark:text-white">
    <svg version="1.1" id="Layer_1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" x="0px" y="0px" width="32px" height="32px" viewBox="0 0 32 32" enable-background="new 0 0 32 32" xml:space="preserve">  <image id="image0" width="32" height="32" x="0" y="0"
    href="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAMAAABEpIrGAAAABGdBTUEAALGPC/xhBQAAACBjSFJN
AAB6JgAAgIQAAPoAAACA6AAAdTAAAOpgAAA6mAAAF3CculE8AAABqlBMVEX////+/v/9/v/0+P35
+/7+///t8/uNt+qbwO2VvOufwu38/f/k7vqDseiWveyGs+kWbNWCsOj6/P78/f74+v7t9Pt6q+cf
ctZ7rOfq8fvu9Pz9/f/+/v5nn+N5que40vJCiN2ewu14quZXleGXvewbb9bZ5/g4gdt9ref9/v5s
ouQccNZ2qObY5vhrouQnd9gxftrQ4faHs+mRuuuPuOsgcteFsulEid0md9h0p+ZqoeQse9kzf9rS
4/c4gtspeNg7hNvU4/cYbdWTu+uSuusjdNdhnOKwzfDa6PgoeNgufNl1qOZpoOQidNc8hNyVvOxR
kd8retk0f9rS4vYXbdWSu+swfdoZbtVmnuM3gdupye82gdtypuVooOR5quY+hdw1gNohc9ePuesV
bNWOuOqUu+u81fKzz/G60/LV5Pc9hNw5g9ve6vlKjd5wpeVto+S30vK81PJBh9yKtekyftrX5vcL
ZdOYvuyixO4WbdVooONUk+Dz9/06g9vd6fg9hdxmn+Olxu7O4Pbh7PnF2vSmx++Queu30fKItOmx
zvD4+v3f6/lkneOkxe4jDKhqAAAAAWJLR0QAiAUdSAAAAAd0SU1FB+cEDgcbF/1gCLEAAAE5SURB
VDjLY2AYBbgAIxAxMeGTZ2bBKckKkmdj52Dg5GJgZGREMRUEuEGGM/Dw8vELCDIIgdVDABMEsTII
A1WKMIiKiUsABSSlYCqkGWRk5UA65BUUGZSUGVRUlYDyauoaQBmIvKaWtg4Dg66evoGhkYCxiakZ
0E41cwtLBiaofitrG1sGdjt7B0d2J2cXVzd3IwYPB08vb7ACJgYZH18HP/+AwCC/YPaQUJdAu7Dw
CI9IhajoGJACRgaZWP24+IREOzs7w6TklFTfyLT0jMys7ByG3DyI6/MFCgpVi4qjS0pVysorKquC
qhlqbGrrGFjqGyAm8DUyNIk1Mwi2MLS2MbR3MHRadXX3qAE9ZdTLwAQPH1BQCcMDDBI0cKY0A6Mw
EwM3EwMTNwOTMCswnliZQCrxRdgoIBkAAJAKLywctTa9AAAAJXRFWHRkYXRlOmNyZWF0ZQAyMDIz
LTA0LTE0VDA3OjI3OjIzKzAwOjAwOxX4lAAAACV0RVh0ZGF0ZTptb2RpZnkAMjAyMy0wNC0xNFQw
NzoyNzoyMyswMDowMEpIQCgAAAAodEVYdGRhdGU6dGltZXN0YW1wADIwMjMtMDQtMTRUMDc6Mjc6
MjMrMDA6MDAdXWH3AAAAAElFTkSuQmCC" />
</svg>

  </span>
</template>
