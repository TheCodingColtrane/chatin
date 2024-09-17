<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch, type PropType, nextTick, h, onBeforeUpdate, onUpdated } from 'vue';
import IconClosedEye from './icons/IconClosedEye.vue';
import IconOpenEye from './icons/IconOpenEye.vue';
import type { RouteLocationNormalizedLoaded } from 'vue-router';
import { backend } from '@/backend/endpoints';
import { useWebSocketStore } from '@/context/websocket';
import MessageEditor from './MessageEditor.vue';
import Modal from './pages/Modal.vue';

const props = defineProps({
  messages: {
    type: Array as PropType<Array<ChatMessages>>,
    required: true
  },
  participants: {
    type: Array as PropType<Array<chatParticipants>>,
    required: true
  },
  user: {
    type: Object as PropType<user>,
    required: true
  },
  chat: {
    type: String,
    required: true
  },
  ws: {
    type: WebSocket,
    required: true
  },
  route: {
    type: Object as PropType<RouteLocationNormalizedLoaded>,
    required: true
  },
  conversation: {
    type: Array as PropType<Array<ChatList>>,
    required: true
  },
})
const messageDetails = reactive(Array<{
  id: string,
  canUpdate: boolean
}>({ id: "", canUpdate: false }))
const participants = reactive({
  senderCode: "",
  receiverCode: [""],
  code: "",
  name: "",
});


const loading = ref(true)
const isMoreMessagesLoaded = ref(false)
const lastMessage = ref('')
const TWENTY_MINUTES = 20 * 60 * 1000;
const contact = ref('')
const websocket = useWebSocketStore().websocket
let recordedAudio = reactive<audio>({ stream: new MediaStream(), recorder: new MediaRecorder(new MediaStream()), blobs: [new Blob()] })
const composedMessage = ref("")
const initalizeChatDetails = () => {
  if (props.messages && props.participants) {
    props.messages.map((c => new Date(c.createdAt).toLocaleTimeString()))
    const messsage = props.messages.find((c) => c.senderCode !== props.user.code)
    if (props.participants.length === 2) {
      participants.receiverCode = [messsage?.senderCode!]
    } else {
      participants.receiverCode = props.participants.filter((c) => c.code !== props.user.code).map((c => c.code))!
    }
    participants.senderCode = props.user.code;
    participants.name = messsage?.fullName!;
    participants.code = messsage?.senderCode!;
    loading.value = false

  }
}

watch(() => props.route.params.code, () => initalizeChatDetails())

const scrollToEnd = () => {
  const chatWindow = document.querySelector('.tab-messages') as HTMLElement;
  chatWindow.scrollTop = chatWindow.scrollHeight;

}


const chatCode = ref(props.route.params.code as string)
onMounted(async () => {
  watch(() => props.messages, (newMessages) => {
    if (newMessages.length > 0 && newMessages[0].code != "") {
      props.messages.map((m => {
        m.fullName = props.participants.find(p => p.code == m.senderCode)?.fullName!
        return m

      }))
      getOlderMessageBatchCode()
      initalizeChatDetails();
      loading.value = false;
    }
  }, { immediate: true });
  await nextTick()
  initializeMessageDetails(props.messages);
  checkMessageUpdateStatus()
  const messageField = document.querySelector('#message-field')! as HTMLElement
  messageField.focus()


})


const initializeMessageDetails = (messages: ChatMessages[]) => {
  messages.forEach((message) => {
    const createdAt = new Date(message.createdAt);
    const currentTime = new Date();
    const elapsedTime = currentTime.getTime() - createdAt.getTime();

    messageDetails.push({ id: message.code, canUpdate: elapsedTime <= TWENTY_MINUTES });
  });
};


const checkMessageUpdateStatus = () => {
  const currentTime = new Date().getTime();
  messageDetails.forEach((detail) => {
    const message = props.messages.find((msg) => msg.code === detail.id);
    if (message) {
      const createdAt = new Date(message.createdAt).getTime();
      detail.canUpdate = (currentTime - createdAt) <= TWENTY_MINUTES;
    }
  });
};

setInterval(checkMessageUpdateStatus, 60 * 1000)

const recordAudioStream = async (recordedAudio: audio): Promise<audio | undefined> => {
  try {
    const stream = await navigator.mediaDevices.getUserMedia({ audio: { echoCancellation: true } });
    const media = new MediaRecorder(stream, { mimeType: 'audio/webm' });
    recordedAudio.stream = stream;
    recordedAudio.recorder = media;
    recordedAudio.blobs = [];

    media.ondataavailable = (evt) => {
      if (evt.data && evt.data.size > 0) {
        recordedAudio.blobs.push(evt.data);
      }
    };

    media.start();
    return recordedAudio;
  } catch (error) {
    console.log(error);
    return undefined;
  }
};

const stopAudioRecordingStream = async (recordedAudio: audio): Promise<audio> => {
  return new Promise((resolve) => {
    recordedAudio.recorder.onstop = () => {
      if (recordedAudio.blobs.length > 0) {
        recordedAudio.stream.getTracks().forEach(track => track.stop());
      }
      resolve(recordedAudio);
    };

    recordedAudio.recorder.stop();
  });
};

const recordAudio = async () => {
  const audio = await recordAudioStream(recordedAudio)!
  if (audio?.recorder) {
    recordedAudio.recorder = audio.recorder
    recordedAudio.stream = audio.stream
  }
}

const stopRecording = async () => {
  if (recordedAudio.recorder.state !== 'inactive') {
    await stopAudioRecordingStream(recordedAudio);
    if (recordedAudio.blobs.length > 0) {
      const blob = new Blob(recordedAudio.blobs, { type: 'audio/webm' });
      sendFileToWS(blob, 'audio/webm');
      recordedAudio.blobs = []
      recordedAudio.recorder = new MediaRecorder(new MediaStream())
      recordedAudio.stream = new MediaStream()
    }
  }
}

const sendFileToWS = (blob: Blob, mimeType: string) => {
  const reader = new FileReader()
  reader.readAsDataURL(blob)
  reader.onloadend = ((evt) => {
    if (reader.result) {
      sendMessage({ mimeType: mimeType, content: reader.result as string, lastModified: new Date(), size: blob.size, name: "", reason: 1 })
    }
  })
}

const sendMessage = (asset?: assets) => {
  const sendMessageButton = document.querySelector('#send-message-button') as HTMLElement
  const action = sendMessageButton.dataset.action!
  const messageField = document.querySelector('#message-field') as HTMLElement
  const messageCode = messageField.dataset.id
  var message = prepareMessage(parseInt(action), messageCode, asset)
  // props.ws.send(JSON.stringify(message))
  websocket?.send(JSON.stringify(message))
  if (asset) {
    props.messages.push({
      code: "", content: "", fullName: props.user.firstName + " " + props.user.lastName,
      createdAt: new Date(), senderCode: props.user.code, asset
    })
  } else {
    props.messages.push({
      code: "", content: composedMessage.value, fullName: props.user.firstName + " " + props.user.lastName,
      createdAt: new Date(), senderCode: props.user.code
    })

  }
  composedMessage.value = ""
  messageField.dataset.id = ""
  sendMessageButton.dataset.action = "0"
}

const prepareMessage = (action: number, code?: string, asset?: assets): message => {
  let message: message
  if (!action) {
    message = {
      type: 0, code: "", content: composedMessage.value,
      senderCode: participants.senderCode, receiverCode: participants.receiverCode, chatCode: chatCode.value /*props.chat*/, seen: false, seenAt: null, asset
    }
  } else {
    message = {
      type: 2, code: code!, content: composedMessage.value,
      senderCode: participants.senderCode, receiverCode: participants.receiverCode, chatCode: chatCode.value /*props.chat*/, seen: true, seenAt: null, asset
    }
  }

  if (asset) {
    if (message.asset?.mimeType === "audio/webm") message.type = 2
    else if (message.asset?.mimeType === "image/jpeg" || message.asset?.mimeType === "image/png") message.type = 3
    else if (message.asset?.mimeType === "video/mp4") message.type = 4
    return message

  }


  return message
}



onMounted(() => {
  document.getElementById('audio-recording-button')?.addEventListener('click', (evt) => {
    const audioButton = evt.target as HTMLElement
    if (audioButton.innerHTML === "RECORD") {
      audioButton.onclick = recordAudio
      audioButton.innerHTML = "STOP RECORDING"

    } else {
      audioButton.onclick = stopRecording
      audioButton.innerHTML = "RECORD"
    }

    nextTick(() => {
      //scrollToEnd()
    }).then()

  })

  document.querySelector('.tab-messages')?.addEventListener('scroll', async (evt) => {
    const currentHeight = evt.target as HTMLElement
    if (lastMessage.value) {
      if (currentHeight.scrollTop == 0) {
        await getMoreMessages(lastMessage.value)
      }
    }

  })


})

onUpdated(() => {
  watch(() => props.route.params.code, (code) => {
    if (code != chatCode.value) {
      chatCode.value = code as string
      isMoreMessagesLoaded.value = false
    }
    getConversationName(chatCode.value)

  }, { immediate: true })
  nextTick(() => {
    getOlderMessageBatchCode()
    if (!isMoreMessagesLoaded.value) scrollToEnd()
  })
})

const getConversationName = (code: string) => {
  const conversation = props.conversation.find(c => c.chatCode === code)
  contact.value = conversation?.groupName || conversation?.recipientFullName || ""
}

const getMoreMessages = async (key: string) => {
  try {
    const messages = await backend.messages.get<ChatMessages[]>(`/${chatCode.value}/messages?k=${key}`)
    if (messages?.length) {
      const lastMsg = document.querySelector(`[data-id=\"${key}\"]`) as HTMLElement
      const rect = lastMsg.getBoundingClientRect()
      let orderedMessages: ChatMessages[] = []
      const messageNumber = messages.length - 1
      for(let i = messageNumber; i >= 0; i--){
        if(i === messageNumber){
          orderedMessages[0] = messages[i]
          continue
        }
        orderedMessages.push(messages[i])
      }

      props.messages.unshift(...orderedMessages)
      const chatWindow = document.querySelector(".tab-messages")!
      let loadedBatches = props.messages.length / 30
      if (loadedBatches % 1 > 0) loadedBatches = loadedBatches - (loadedBatches % 1)
      await nextTick()
      chatWindow.scrollTop = chatWindow.scrollHeight / loadedBatches - rect.y
      getOlderMessageBatchCode()
      isMoreMessagesLoaded.value = true
    }
  } catch (error) {

  }

}

const getFile = (event: Event) => {
  const selection = event.target as HTMLInputElement
  const selectedFiles = selection.files
  if (selectedFiles) {
    const file = selectedFiles.item(0)
    sendFileToWS(file!, selectedFiles.item(0)?.type!)
  }
  return
}

document.querySelector('#upload-files')?.addEventListener('change', (e) => {
  const selection = e.target as HTMLInputElement
  const selectedFiles = selection.files
  if (selectedFiles) {
    const file = selectedFiles.item(0)
    sendFileToWS(file!, selectedFiles.item(0)?.type!)
  }
})

document.querySelector('#message-field')?.addEventListener('keypress', (e) => {
  const messageField = e.target as HTMLInputElement
  if (!messageField.value.length) {
    messageField.dataset.id = ""
  }
})


const formatTime = (time: Date) => {
  return new Date(time).toLocaleTimeString();
}

const setMessageUpdate = (content: string, code: string) => {
  composedMessage.value = content
  const messageField = document.querySelector('#message-field') as HTMLElement
  const sendMessageButton = document.querySelector('#send-message-button') as HTMLElement
  sendMessageButton.dataset.action = "1"
  messageField.dataset.id = code

}

const getOlderMessageBatchCode = () => {
  const messageCount = props.messages.length
  if (messageCount >= 30) {
    const remainingMessages = messageCount % 30
    if (remainingMessages === 0) {
      const key = props.messages[0].code
      lastMessage.value = key
    }
  }
}


const canUpdate = (code: string) => messageDetails.find(md => md.id === code)?.canUpdate ?? false

</script>

<template>

  <div class="chat-window" v-if="messageDetails">
    <h2 class="conversation-name">{{ contact }}</h2>
    <div class="messages">

    <div v-for="message in props.messages">
      <div class="container darker" v-if="message.senderCode === props.user.code">
        <article v-if="message.content" :data-id="message.code" :data-seen="message.seen">{{ message.content }}
        </article>
        <div v-if="message.asset">
          <div v-if="message.asset.mimeType === 'audio/webm'">
            <audio controls :src=message.asset?.name></audio>
          </div>
          <div v-else-if="message.asset.mimeType === 'image/jpeg' || message.asset.mimeType === 'image/png'">
            <img :src="message.asset?.name">
          </div>
          <div v-else-if="message.asset.mimeType === 'video/mp4'">
            <video :src="message.asset.name"></video>
          </div>
        </div>
        <div>
          <span class="time-left">{{ formatTime(message.createdAt) }}</span>
          <button type="button" v-if="canUpdate(message.code)"
            @click="setMessageUpdate(message.content, message.code)">Update Message</button>
          <IconOpenEye v-if="message.seen"></IconOpenEye>
          <IconClosedEye v-else></IconClosedEye>
          <span v-if="message.updatedAt">MENSAGEM EDITADA EM : {{ message.updatedAt }}</span>
        </div>
      </div>
      <div class="container" v-else>
        <article v-if="message.content" :data-id="message.code" :data-seen="message.seen">{{ message.content }}
        </article>
        <div v-if="message.asset?.name">
          <div v-if="message.asset.mimeType === 'audio/webm'">
            <audio controls :src=message.asset?.name></audio>
          </div>
          <div v-else-if="message.asset.mimeType === 'image/jpeg' || message.asset.mimeType === 'image/png'">
            <img :src="message.asset?.name">
          </div>
          <div v-else-if="message.asset.mimeType === 'video/mp4'">
            <video :src="message.asset.name"></video>
          </div>
        </div>
        <span class="time-right">{{ formatTime(message.createdAt) }}</span>
        <button type="button" v-if="canUpdate(message.code)"
          @click="setMessageUpdate(message.content, message.code)">Update Message</button>
        <IconOpenEye v-if="message.seen"></IconOpenEye>
        <IconClosedEye v-else></IconClosedEye>
        <span v-if="message.updatedAt">MENSAGEM EDITADA EM : {{ new Date(message.updatedAt).toLocaleString() }}</span>
      </div>
    </div>
    <div class="input-container">


      <div class="message-input">
        <!-- <MessageEditor v-model="composedMessage" /> -->

        <!-- <Editor>

          </Editor> -->
          <input type="text" v-model="composedMessage" id="message-field" data-id="" placeholder="Type your message..." style="width: 100%; height: 50px;"/>
          <!-- <input type="file" id="upload-files" @change="getFile">
          <button @click="sendMessage()" id="send-message-button" data-action="0">Send</button>
          <button @click="recordAudio" id="audio-recording">RECORD</button>  -->
      </div>
      <div class="action-buttons">
        <input type="file" id="upload-files" @change="getFile">
        <button @click="sendMessage()" id="send-message-button" data-action="0">Send</button>
        <button @click="recordAudio" id="audio-recording-button">RECORD</button>
      </div>
    </div>
    </div>
  </div>
</template>



<style scoped>
.action-buttons {
  display: flex;
  justify-content: flex-start;
}

.chat-window {
  display: flex;
  flex-direction: column;
  /* height: 100vh; */
  overflow: hidden;
  /* Prevents the chat window from overflowing */
}

.conversation-name {
  position: fixed;
  top: 0;
  left: 30%;
  width: 70%;
  padding: 10px 12px;
  background-color: white;
  border-bottom: 1px solid #ccc;
  z-index: 1;
  height: 40px;
  display: flex;
  align-items: center;
}

.messages {
  flex: 1;
  overflow-y: auto;
  padding-top: 60px; /* Adjust based on the height of the conversation name */
  padding-bottom: 120px;
}

.input-container {
  position: fixed;
  bottom: 0;
  left: 30%;
  width: 70%;
  height: 120px;
  background-color: white;
  display: flex;
  flex-direction: column;
  padding: 10px;
  box-shadow: 0 -1px 5px rgba(0, 0, 0, 0.1);
  z-index: 1;
}

.message-input {
  flex: 1;
  margin-bottom: 10px;
  height: 70px;
}

#send-message-button{
  padding: 5px 10px;
  border: none;
  border-radius: 3px;
  background-color: #333;
  color: #fff;
  cursor: pointer;
}

#audio-recording-button {
  padding: 5px 10px;
  border: none;
  border-radius: 3px;
  background-color: #333;
  color: #fff;
  cursor: pointer;
  margin-left: 5px;
}

.container {
  border: 2px solid #dedede;
  background-color: #f1f1f1;
  border-radius: 5px;
  padding: 10px;
  margin: 15px 0;
}

.darker {
  border-color: #ccc;
  background-color: #ddd;
  margin: 15px;
}

.container::after {
  content: "";
  clear: both;
  display: table;
}

.container img {
  float: left;
  max-width: 60px;
  width: 100%;
  margin-right: 20px;
  border-radius: 50%;
}

.container img.right {
  float: right;
  margin-left: 20px;
  margin-right: 0;
}

.time-right {
  float: right;
  color: #aaa;
}

.time-left {
  float: left;
  color: #999;
}

.tab-messages>.chat-window {
  overflow-y: auto;
  flex-grow: 1;
  margin-bottom: 10px;
}




</style>