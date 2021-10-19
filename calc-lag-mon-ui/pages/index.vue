<template>
    <h1>Hello outer</h1>
</template>

<script lang="ts">
export default defineComponent({
    mounted() {
        if (!this.$route.fullPath.startsWith('/clm-ui'))
            this.$router.push('/clm-ui');
    }
});
</script>