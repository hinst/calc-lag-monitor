<template>
  <div>
    calculation lag chart
  </div>
</template>

<script lang="ts">
import { defineComponent, useContext } from '@nuxtjs/composition-api'

export default defineComponent({
  setup() {
    const context = useContext();
    setTimeout(async () => {
      const response = await fetch(context.env.apiUrl + '/clm/lag');
      console.log(response.statusText);
    }, 2000);
  },
});
</script>
