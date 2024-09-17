import backend from '@/backend';
import { defineStore } from 'pinia';
import { ref, onUnmounted } from 'vue';


export const useWebSocketStore = defineStore('websocket', () => {
  const websocket = ref<WebSocket | null>(null);
  const messages = ref<ChatMessages[]>([]);
  const chatList = ref<ChatList[]>([]);
  const currentUser = ref<user | null>(null);
  const loading = ref(true);
  const connectionStatus = ref<'connecting' | 'connected' | 'disconnected' | 'error'>('connecting');
  const errorMessage = ref<string | null>(null);

  const connectWebSocket = (url: string) => {
    try {
      websocket.value = new WebSocket(!url ? 'ws://localhost:8080/chat' : url);

      websocket.value.addEventListener('open', () => {
        console.log('User just connected to chat');
        connectionStatus.value = 'connected';
        loading.value = false;
      });

      websocket.value.addEventListener('close', () => {
        connectionStatus.value = 'disconnected';
      });

      websocket.value.addEventListener('error', (evt) => {
        console.error('An error occurred while using the chat:', evt);
        connectionStatus.value = 'error';
        errorMessage.value = 'An error occurred. Please try again later.';
      });

      websocket.value.addEventListener('message',    (evt) => {
        handleWebSocketMessage(evt);
      });

    } catch (error) {
      console.error('Failed to connect to WebSocket:', error);
      connectionStatus.value = 'error';
      errorMessage.value = 'Failed to connect to WebSocket. Please check your connection.';
      loading.value = false;
    }
  };

  const disconnectWebSocket = () => {
    if (websocket.value) {
      websocket.value.close();
      websocket.value = null;
      connectionStatus.value = 'disconnected';
    }
  };

  const handleWebSocketMessage = (evt: MessageEvent) => {
    const receivedMessage = JSON.parse(evt.data) as WebSocketResponse;
    const user = messages.value.find((c) => c.senderCode !== currentUser.value?.code);
    if (user) {
      receivedMessage.fullName = user.fullName;
      receivedMessage.userCode = user.senderCode;
    }

    switch (receivedMessage.type) {
      case 0:
        createMessage(receivedMessage);
        break;
      case 1:
        uptadeMessageSeenStatus(receivedMessage);
        break;
      case 2:
        uodateMessageContent(receivedMessage);
        break;
      case 4:
        // Handle type 4 if necessary
        break;
      default:
        console.warn('Unknown message type:', receivedMessage.type);
    }
  };

  const createMessage  = (receivedMessage: WebSocketResponse) => {
    messages.value.push({
      code: receivedMessage.code,
      content: receivedMessage.content,
      createdAt: receivedMessage.createdAt,
      fullName: receivedMessage.fullName!,
      senderCode: receivedMessage.userCode!,
    });

    websocket.value?.send(
      JSON.stringify({
        code: receivedMessage.messageCode,
        type: 5,
        receiverCode: messages.value.map((c) => c.senderCode),
      })
    );

    chatList.value.push({
      messageContent: receivedMessage.content,
      userCode: receivedMessage.userCode!,
      recipientCode: receivedMessage.userCode!,
      recipientFullName: receivedMessage.fullName!,
      messageSeen: false,
      senderCode: currentUser.value?.code!,
      senderFullName: `${currentUser.value?.firstName} ${currentUser.value?.lastName}`,
      messageCreatedAt: receivedMessage.createdAt,
      assetName: receivedMessage.asset?.name ?? '',
      assetMimeType: receivedMessage.asset?.mimeType ?? '',
      groupName: '',
      groupCode: '',
      chatCode: receivedMessage.code,
    });
  };

  const uptadeMessageSeenStatus = (receivedMessage: WebSocketResponse) => {
    const message = messages.value.find((c) => c.code === receivedMessage.code);
    if (message) message.seen = true;
  };

  const uodateMessageContent = (receivedMessage: WebSocketResponse) => {
    const message = messages.value.find((c) => c.code === receivedMessage.code);
    if (message) message.content = receivedMessage.content;
  };

  onUnmounted(() => {
    disconnectWebSocket();
  });

  return {
    websocket,
    messages,
    chatList,
    currentUser,
    loading,
    connectionStatus,
    errorMessage,
    connectWebSocket,
    disconnectWebSocket,
  };
});
