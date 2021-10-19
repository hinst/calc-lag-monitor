<template>
    <h1>Hello outer</h1>
</template>

<script lang="ts">
export default defineComponent({
    mounted() {
        console.log('mounted');
    }
});
</script>