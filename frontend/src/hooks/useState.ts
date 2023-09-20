import { ref, type Ref } from 'vue'

export const useState = (initState: boolean): [Ref<boolean>, (bool: boolean) => void] => {
  const state = ref(initState)
  const updateState = (_state: boolean) => {
    state.value = _state
  }
  return [state, updateState]
}
