<script lang="ts" setup>
import { useGettext } from 'vue3-gettext'

const { $gettext } = useGettext()
const year = new Date().getFullYear()
</script>

<template>
  <footer class="text-sm text-text-tertiary f-c-c flex-col gap-1">
    <p>
      © 2022 - {{ year }}
      <a
        class="hover:text-brand-hover hover:underline"
        rel="noreferrer"
        target="__blank"
        href="https://ornaverse.com/"
      >
        {{ $gettext('OrnaVerse Panel') }}
      </a>
      {{ $gettext('All Rights Reserved.') }}
    </p>
  </footer>
</template>
