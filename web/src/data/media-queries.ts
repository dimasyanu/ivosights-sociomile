import { useMediaQuery } from "@vueuse/core"
import { computed } from "vue"

export const _xs = useMediaQuery("screen and (min-width: 0rem)")
export const _sm = useMediaQuery("screen and (min-width: 40rem)")
export const _md = useMediaQuery("screen and (min-width: 48rem)")
export const _lg = useMediaQuery("screen and (min-width: 64rem)")
export const _xl = useMediaQuery("screen and (min-width: 80rem)")
export const _x2l = useMediaQuery("screen and (min-width: 96rem)")

export const xs = computed(() => _xs.value && !_sm.value)
export const sm = computed(() => _sm.value && !_md.value)
export const md = computed(() => _md.value && !_lg.value)
export const lg = computed(() => _lg.value && !_xl.value)
export const xl = computed(() => _xl.value && !_x2l.value)
export const x2l = computed(() => _x2l.value)
