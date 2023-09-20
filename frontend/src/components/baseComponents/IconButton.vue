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
.btn {
  padding: 0 6px;
  min-width: 28px;
  min-height: 28px;
  text-align: center;
  font-size: 14px;
  color: #8d8e9a;
  background: #101217;
  border: 1px solid #22252a;
  border-radius: 1px;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  user-select: none;
  transition: all 0.2s;

  &:hover {
    opacity: 0.6;
  }

  .text {
    padding: 0 4px 0 8px;
    font-size: 12px;
  }
}
.slot {
  position: absolute;
  background: #101217;
  left: 0;
  right: 0;
  margin: 0 6px;
}
</style>
