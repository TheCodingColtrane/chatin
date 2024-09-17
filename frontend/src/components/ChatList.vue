<script lang="ts" setup>
import {  onMounted , ref, watch, type PropType } from 'vue';
import { useRouter, type RouteLocationNormalizedLoaded } from 'vue-router';
import Chat from './Chat.vue';
import { useWebSocketStore } from '@/context/websocket';

const props = defineProps({
  conversations: {
    type: Array as PropType<Array<ChatList>>,
    required: true
  },
  route: {
    type: Object as PropType<RouteLocationNormalizedLoaded>,
    required: true
  },
  user: {
    type: Object as PropType<user>,
    required: true
  },
  messages: {
    type: Array as PropType<Array<ChatMessages>>,
    required: true
  },
  participants: {
    type: Array as PropType<Array<chatParticipants>>,
    required: true
  },
  ws: {
    type: WebSocket,
    required: true
  }
})

const contact = ref('')
const search = ref('')
const router = useRouter()
const websocket = useWebSocketStore().websocket
const openConversation = (code: string) => {
  router.replace({ name: "chat", params: { code } })
  props.route.params.id = code
}

const getContactName = (chatCode: string) => {
  const conversation = props.conversations.find(c => c.chatCode === chatCode)
  const groupName = conversation?.groupName
  const recipient = conversation?.recipientFullName!
  contact.value = groupName || recipient
  return groupName || recipient 
}

onMounted(() => {
  watch(() => props.route.params.code, (code) => {
  const conversationsList = document.getElementsByClassName("conversation-item");
  for (let i = 0; i < conversationsList.length; i++) {
    if(conversationsList[i].className === "conversation-item active"){
    conversationsList[i].className = conversationsList[i].className.replace("active", "");
      const activeConversation = conversationsList[i] as HTMLElement
      activeConversation.style.display = "block"
    }
  }
  const conversationItem = document.querySelector('[data-id=' + code + ']')!
  conversationItem.className = "conversation-item active"
  getContactName(code as string)
} )
})




</script>

<template>
  <div class="tab">
    <input type="search" v-model="search" id="message-field" data-id="" placeholder="Search" style="width: 100%;"/>
    <div v-for="conversation in conversations!">
      <li @click="openConversation(conversation.chatCode)" :data-id="conversation.chatCode" class="conversation-item">
        <h4>{{ getContactName(conversation.chatCode) }}</h4>
        <span>{{ conversation.senderCode === props.user.code ? "Você: " + conversation.messageContent  : conversation.senderFullName + ": " +
          conversation.messageContent }}</span>
      </li>
    </div>
  </div>
  <div class="tab-messages">
    <Chat :messages="props.messages" :participants="props.participants" :user="user"
      :chat="route.params.code.toString()" :route="route" :ws="props.ws" :conversation="conversations" v-if="contact"></Chat>
  </div>

</template>


<style>
* {
  box-sizing: border-box
}

.tab {
  float: left;
  border: 1px solid #ccc;
  background-color: #f1f1f1;
  width: 30%;
  height: 100vh; 
  /* overflow-y: auto;  */
  position: fixed; 
}

.tab li {
  display: block;
  background-color: inherit;
  color: black;
  padding: 5.5px 4px;
  width: 100%;
  border: none;
  outline: #ddd 1px solid;
  text-align: left;
  cursor: pointer;
  transition: 0.3s;
  font-size: 17px;
}

.tab li:hover {
  background-color: #ddd;
}

.tab li.active {
  background-color: #ccc;
}

.tab-messages {
   margin-left: 30%; /* Move the content to the right, next to the fixed menu */
  padding: 0px 12px;
  padding-top: 55px; /* Adjust for the height of the fixed contact name */
  border: 1px solid #ccc;
  width: 70%;
  border-left: none;
  height: calc(100vh - 55px); /* Adjust height to account for the fixed contact name */
  overflow-y: scroll
}
</style>