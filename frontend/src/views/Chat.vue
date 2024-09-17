<script setup lang="ts">
import { backend } from '@/backend/endpoints';
import { onMounted, ref, reactive, watch } from 'vue';
import { useRoute } from 'vue-router';
import ChatList from '@/components/ChatList.vue';
import { useWebSocketStore } from '@/context/websocket';
const route = useRoute()
const chatMessages = reactive<ChatDetails>({
  messages: [{
    code: "", content: "", createdAt: new Date(), fullName: "", senderCode: "", seen: false,
    asset: { mimeType: "", name: "", reason: 0 }
  }], participants: [{ code: "", fullName: "", profileImg: "", username: "" }]
});
let chatList = reactive<ChatList[]>([{
  assetMimeType: "", assetName: "", chatCode: "", groupCode: "", groupName: "", messageContent: "", messageCreatedAt: new Date,
  messageSeen: false, userCode: "", senderCode: "", senderFullName: "", recipientCode: "", recipientFullName: ""
}])
let loading = ref(true)
const websocketStore = useWebSocketStore()
websocketStore.connectWebSocket("")
let websocket = websocketStore.websocket




let currentUser = reactive<user>({
  code: "",
  email: "",
  firstName: "",
  lastName: "",
  username: "",
  createdAt: new Date(),
});

watch(
  () => route.params.code,
  async (code) => {
    chatMessages.messages = []
    chatMessages.participants = []
    await getChatDetails(code as string)
  }


)


const getChatList = async (): Promise<ChatList[]> => {
  const chats = await backend.messages.get<ChatList[]>("/")
  return chats ?? []
}

const getChatDetails = async (code: string) => {
  try {
    const [chats, userData] = await Promise.all([
      backend.messages.get<ChatDetails>("/" + code),
      backend.users.get<user>("/data"),
    ]);

    if (chats && userData) {
      const { participants } = chats
      const foundChats = chats.messages.map((c, i) => {
        if (c.asset?.name) {
          c.asset.name = "http://localhost:8001/" + c.asset?.name
        }
        return c
      })
      chatMessages.messages = foundChats.reverse();
      chatMessages.participants = participants

      currentUser = userData
    }
  } catch (error) {
    console.error("Error fetching data:", error);
    loading.value = false
  }
}


onMounted(async () => {
  try {
    chatList = await getChatList()
    if (route.params.code) {
      const code = route.params.code as string
      await getChatDetails(code)
    }

    const webscoket = useWebSocketStore().connectWebSocket('ws://localhost:8080/chat')
    // websocket.value.addEventListener('open', (evt) => {
    //   console.log('User just connected to chat')
    //   // seeMessage()
    // })

    // websocket.value.addEventListener('close', (evt) => {
    // })

    // websocket.value.addEventListener('error', (evt) => {

    //   console.log('An error had ocurred while using the chat')
    // })

    // websocket.value.addEventListener('message', async (evt) => {
    //   const receivedMessage = JSON.parse(evt.data) as WebSocketResponse

    //   const user = chatMessages.messages.find((c) => c.senderCode !== currentUser.code)
    //   receivedMessage.fullName = user?.fullName!
    //   receivedMessage.userCode = user?.senderCode!
    //   if (receivedMessage.type === 0) {
    //     chatMessages.messages.push({
    //       code: receivedMessage.code, content: receivedMessage.content, createdAt: receivedMessage.createdAt,
    //       fullName: receivedMessage.fullName, senderCode: receivedMessage.userCode
    //     })
    //     websocket.value.send(JSON.stringify({ code: receivedMessage.messageCode, type: 5, receiverCode: chatMessages.participants.map(c => c.code) }))
    //     chatList.push({
    //       messageContent: receivedMessage.content, userCode: receivedMessage.userCode, recipientCode: receivedMessage.userCode, recipientFullName: receivedMessage.fullName,
    //       messageSeen: false, senderCode: currentUser.code, senderFullName: currentUser.firstName + " " + currentUser.lastName, messageCreatedAt: receivedMessage.createdAt,
    //       assetName: receivedMessage.asset?.name ?? "", assetMimeType: receivedMessage.asset?.mimeType ?? "",
    //       groupName: "", groupCode: "", chatCode: receivedMessage.code
    //     })

    //   }

    //   else if (receivedMessage.type === 1) {
    //     chatMessages.messages.find(c => c.code === receivedMessage.code)!.seen = true
    //     const lastMessageIndex = chatMessages.messages.length - 1
    //     chatMessages.messages[lastMessageIndex].code = receivedMessage.code
    //   } else if (receivedMessage.type === 2) {
    //     chatMessages.messages.find(c => c.code === receivedMessage.code)!.content = receivedMessage.content
    //   } else if (receivedMessage.type === 4) {

    //   }
    // })

    loading.value = false

  }
  catch (error) {
    console.error("Error fetching data:", error);
    loading.value = false
  }
})



</script>

<template>
  <div v-if="!loading && websocket?.readyState === websocket?.OPEN && chatMessages">
    <ChatList :conversations="chatList" :route="route" :user="currentUser" :messages="chatMessages.messages"
      :participants="chatMessages.participants" :ws="websocket!" />
    <!-- <div v-if="!loading && websocket.readyState === websocket.OPEN && chatMessages"> -->
    {{ console.log("olha o dados", chatMessages.messages.length) }}
    <!-- <Chat :messages="chatMessages.messages" :participants="chatMessages.participants" :user="currentUser" :chat="code"
      :route="route" :ws="websocket"/> -->
  </div>

</template>

<style>
html, body {
  margin: 0;
  padding: 0;
  height: 100%;
  overflow: hidden; 
}
</style>