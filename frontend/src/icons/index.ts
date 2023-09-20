import CloseIcon from './CloseIcon.vue'
import DropdownIcon from './DropdownIcon.vue'
import OpenIcon from './OpenIcon.vue'
import RefreshIcon from './RefreshIcon.vue'
import SyncIcon from './SyncIcon.vue'
import GithubIcon from './GithubIcon.vue'

export const Icons = {
  close: CloseIcon,
  dropdown: DropdownIcon,
  open: OpenIcon,
  refresh: RefreshIcon,
  sync: SyncIcon,
  github: GithubIcon
}

export type IconsKey = keyof typeof Icons
