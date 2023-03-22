<script setup lang="ts">
defineProps<{
  label: string
  selected: boolean
}>()

defineEmits<{
  click: () => void
}>()
</script>

<template>
  <div
    :class="{
      'tab-item': true,
      'tab-item--selected': selected
    }"
    @click="$emit('click')"
  >
    <i class="tab-item__icon">
      <slot />
    </i>
    <div class="tab-item__label">{{ label }}</div>
  </div>
</template>

<style scoped>
.tab-item {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: flex-start;
  padding: 0 1.25rem 0 1rem;
  height: 100%;
  cursor: pointer;
}

.tab-item:hover {
  background-color: var(--color-background-soft);
}

.tab-item__icon {
  width: 1.5rem;
  height: 1.5rem;
  color: var(--color-text);
}

.tab-item__label {
  font-size: 1rem;
  margin-left: 0.75rem;
}

.tab-item--selected {
  box-shadow: 0 -2px var(--color-primary) inset;
}

.tab-item--selected .tab-item__icon {
  color: var(--color-primary);
}

.tab-item--selected .tab-item__label {
  color: var(--color-primary);
}
</style>
