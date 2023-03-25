<script setup lang="ts">
import NavigatePreviousIcon from './icons/IconNavigatePrevious.vue'
import NavigateNextIcon from './icons/IconNavigateNext.vue'

defineProps<{
  isNextDisabled: boolean
  isPreviousDisabled: boolean
  label: string
}>()

defineEmits<{
  previous: () => void
  next: () => void
}>()
</script>

<template>
  <div class="navigation">
    <div class="navigation__label">{{ label }}</div>
    <i
      class="navigation__previous"
      :class="{
        'navigation__previous--disabled': isPreviousDisabled
      }"
      @click="() => $emit('previous')"
    >
      <NavigatePreviousIcon />
    </i>
    <i
      class="navigation__next"
      :class="{
        'navigation__next--disabled': isNextDisabled
      }"
      @click="() => $emit('next')"
    >
      <NavigateNextIcon />
    </i>
  </div>
</template>

<style scoped>
.navigation {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
}

.navigation__label {
  font-family: 'IBM Plex Mono', monospace;
  font-size: 0.8rem;
}

i {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 1.5rem;
  height: 1.5rem;
  cursor: pointer;
}

i:hover {
  color: var(--color-primary);
}

.navigation__previous--disabled,
.navigation__next--disabled,
i:hover.navigation__previous--disabled,
i:hover.navigation__next--disabled {
  color: var(--color-background-mute);
  cursor: not-allowed;
}
</style>
