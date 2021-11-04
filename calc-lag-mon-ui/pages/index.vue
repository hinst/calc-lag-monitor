<template>
  <h1>Redirecting...</h1>
</template>

<script lang="ts">
import { defineComponent, useRouter } from '@nuxtjs/composition-api';
export default defineComponent({
  setup() {
    const router = useRouter();
    router.push('/calculation-lag-chart');
    return {};
  }
});
</script>
