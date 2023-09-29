<script setup lang="ts">
import { ref } from 'vue'
import { Icons, type IconsKey } from '@/icons'

interface Props {
  icon: IconsKey
  text?: string
  type?: 'button' | 'link'
  href?: string
}

withDefaults(defineProps<Props>(), {
  type: 'button'
})

const innerShowSlot = ref(false)

const handleClick = () => {
  innerShowSlot.value = true
}

const onMouseLeave = () => {
  innerShowSlot.value = false
}
</script>

<template>
  <div @mouseleave="onMouseLeave" style="position: relative">
    <template v-if="type === 'link'">
      <a :href="href" target="_blank" class="btn"> <component :is="Icons[icon]" /></a>
    </template>
    <template v-else-if="type === 'button'">
      <div @click="handleClick" class="btn">
        <div v-if="text" class="text">{{ text }}</div>
        <component :is="Icons[icon]" />
      </div>
    </template>

    <div v-show="innerShowSlot" class="slot">
      <slot></slot>
    </div>
  </div>
</template>

<style lang="less" scoped>
@import url(style.less);
</style>
