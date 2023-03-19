<script setup lang="ts">
defineProps<{
  label: string
  selected: boolean
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
  padding: 0 16px;
  height: 100%;
  cursor: pointer;
}

.tab-item:hover {
  background-color: var(--color-background-soft);
}

.tab-item__icon {
  width: 24px;
  height: 24px;
  color: var(--color-text);
}

.tab-item__label {
  font-size: 1rem;
  margin-left: 12px;
}

.tab-item--selected {
  box-shadow: 0px -2px var(--color-primary) inset;
}

.tab-item--selected .tab-item__icon {
  color: var(--color-primary);
}

.tab-item--selected .tab-item__label {
  color: var(--color-primary);
}

@media (min-width: 1024px) {
  .tab-item {
    padding: 0 24px;
  }
}
</style>
