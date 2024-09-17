<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import { backend } from '@/backend/endpoints';
import { createCookie } from '@/utils/cookies';
import Modal from '@/components/pages/Modal.vue';

const loginData = reactive({ email: "", password: "" })
const login = async () => {
    const authorization = await backend.auth.post(loginData) as any
    if (authorization) {
        createCookie("Authorization", authorization.token, authorization.exp)

        window.location.href = "/chat"
    }
}
</script>

<template>
    <form>
        <label>E-mail</label>
        <input type="email" v-model="loginData.email" id="">
        <label>Password</label>
        <input type="email" v-model="loginData.password" id="">
        <input type="button" @click="login">
    </form>
</template>
<style scoped>
form {
    border: 1px solid black;
    border-radius: border-box;
    z-index: 9999999;
    width: 300px;
}
</style>