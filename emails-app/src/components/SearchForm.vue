<script setup lang="ts">
import { ref } from 'vue'

import type { Query } from '@/models/query'

const from = ref('')
const to = ref('')
const subject = ref('')
const bodyIncludes = ref('')
const bodyExcludes = ref('')
const dateFrom = ref()
const dateTo = ref()

const buildQuery = () => {
  return {
    from: from.value !== '' ? from.value : undefined,
    to: to.value !== '' ? to.value.split(',').map((s) => s.trim()) : undefined,
    subjectIncludes: subject.value !== '' ? subject.value : undefined,
    bodyIncludes: bodyIncludes.value !== '' ? bodyIncludes.value : undefined,
    bodyExcludes: bodyExcludes.value !== '' ? bodyExcludes.value : undefined,
    dateRange:
      dateFrom.value !== '' || dateTo.value !== ''
        ? {
            from: dateFrom.value !== '' ? new Date(dateFrom.value) : undefined,
            to: dateTo.value !== '' ? new Date(dateTo.value) : undefined
          }
        : undefined
  } as Query
}

const reset = () => {
  from.value = ''
  to.value = ''
  subject.value = ''
  bodyIncludes.value = ''
  bodyExcludes.value = ''
  dateFrom.value = ''
  dateTo.value = ''
}
</script>

<template>
  <!-- The search form component renders a form with input fields for the query parameters. 
  It emits a search event when the user clicks the search button.
  It resets the form when the user clicks the reset button.
  The form is not submitted when the user presses the enter key.
  It's only supposed to build a query object and emit it. -->

  <div>
    <form @submit.prevent>
      <div class="form-group">
        <label for="from">From</label>
        <input type="text" id="from" v-model="from" />
      </div>

      <div class="form-group">
        <label for="to">To</label>
        <input type="text" id="to" v-model="to" />
      </div>

      <div class="form-group">
        <label for="subject">Subject</label>
        <input type="text" id="subject" v-model="subject" />
      </div>

      <div class="form-group">
        <label for="body-includes">Body includes</label>
        <input type="text" id="body-includes" v-model="bodyIncludes" />
      </div>

      <div class="form-group">
        <label for="body-excludes">Body excludes</label>
        <input type="text" id="body-excludes" v-model="bodyExcludes" />
      </div>

      <div class="form-group form-group--date">
        <label for="date-from">Date from</label>
        <input type="date" id="date-from" v-model="dateFrom" />
      </div>

      <div class="form-group form-group--date">
        <label for="date-to">Date to</label>
        <input type="date" id="date-to" v-model="dateTo" />
      </div>

      <div class="form-group form-group--buttons">
        <button class="form-group--buttons__reset" type="button" @click="reset">Reset</button>
        <button class="form-group--buttons__search" type="button" @click="buildQuery">
          Search
        </button>
      </div>
    </form>
  </div>
</template>

<style scoped>
.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.form-group label {
  font-size: 0.8rem;
  font-weight: 500;
  color: var(--color-heading);
}

.form-group input {
  padding: 0.5rem;
  border: 1px solid var(--color-border);
  border-radius: 0.5rem;
  outline: none;
  font-size: 1rem;
  font-weight: 500;
  color: var(--color-heading);
}

.form-group--date {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.form-group--date label {
  font-size: 0.8rem;
  font-weight: 500;
  color: var(--color-heading);
}

.form-group--date input {
  padding: 0.5rem;
  border: 1px solid var(--color-border);
  border-radius: 0.5rem;
  outline: none;
  font-size: 1rem;
  font-weight: 500;
  color: var(--color-heading);
}

.form-group--buttons {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
}

.form-group--buttons button {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 0.5rem;
  outline: none;
  font-size: 1rem;
  font-weight: 500;
  color: var(--color-heading);
  background-color: var(--color-primary);
  cursor: pointer;
}

.form-group--buttons button:hover {
  background-color: var(--color-primary-hover);
}

.form-group--buttons__reset {
  background-color: var(--color-secondary);
}

.form-group--buttons__reset:hover {
  background-color: var(--color-secondary-hover);
}

.form-group--buttons__search {
  background-color: var(--color-primary);
}

.form-group--buttons__search:hover {
  background-color: var(--color-primary-hover);
}

.form-group--buttons__search:disabled {
  background-color: var(--color-primary-disabled);
  cursor: not-allowed;
}

.form-group--buttons__search:disabled:hover {
  background-color: var(--color-primary-disabled);
}
</style>
