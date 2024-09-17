import backend from '@/backend';
import { defineStore } from 'pinia';
import { ref } from 'vue';
import {useRouter} from 'vue-router'

const router = useRouter()

export const useUserStore = defineStore('user', () => {
    const user = ref<user | null>(null)
    const $reset = () => {
        user.value = null
    }

    const get = async(): Promise<user | null> => {
        try {   
            const data = await backend.get<user>('/user/data');
            if(!data){
                router.push('/login')
            }
            return data!
            
        } catch (error) {
            return null
        }

    }

    return {
        user, 
        get,
        $reset
    }

    

    
})
