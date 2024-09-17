<script lang="ts" setup>
const props = defineProps({
    message: {
        type: String,
        required: true
    },
    title: {
        type: String,
        required: true,
    },
    timer: {
        type: Number,
        required: false
    }
})


</script>

<template>
    <div class="error-message">
        <div class="error-title">{{ props.title }}</div>
        <p>{{props.message}}</p>
    </div>
</template>


<style scoped>
.error-block {
    margin: 28px 0;
    padding: 20px 24px 4px 42px;
    border-radius: 8px;
    overflow-x: auto;
    transition: color .5s, background-color .5s;
    position: relative;
    font-size: 14px;
    line-height: 1.6;
    font-weight: 500;
    color: #0000008c;
    background-color: #ed3c50;
}

.error-title {
    margin-bottom: 8px;
    font-size: 15px;
    font-weight: 500;
    color: #ed3c50;
    transition: color .5s;
}
</style>