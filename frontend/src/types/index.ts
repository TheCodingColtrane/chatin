type message = {
  type: number
  code: string
  content: string
  senderCode: string
  seen: boolean | null
  seenAt: Date | null;
  // receiverCode: string
  receiverCode: string[]
  chatCode: string
  asset?: assets
}

type audio = {
  stream: MediaStream
  recorder: MediaRecorder
  blobs: Blob[]
}


type assets = {
  mimeType: string
  name: string
  reason: number;
  size?: number
  content?: string;
  lastModified?: Date,
}

type chatAssets = {
  mimeType: string
  name: string,
  reason: number
}

type chatParticipants = {
  code: string
  fullName: string;
  username: string;
  profileImg: string;
}

type chatMessages = {
  code: string
  messageCode?: string;
  userCode: string
  username: string
  fullName: string
  content: string
  seen: boolean
  seenAt?: boolean
  createdAt: Date
  updatedAt?: Date
  asset?: chatAssets
}

interface WebSocketResponse {
  code: string
  messageCode?: string;
  userCode: string
  username: string
  fullName: string
  content: string
  seen: boolean
  seenAt?: boolean
  createdAt: Date
  updatedAt?: Date
  asset?: chatAssets
  type: number;
}


type chatDetails = {
  messages: Array<chatMessages>,
  participants: Array<chatParticipants>
}

interface ChatMessages {
  code: string
  messageCode?: string
  senderCode: string
  fullName: string
  seen?: boolean
  seenAt?: Date
  content: string
  createdAt: Date
  updatedAt?: Date
  asset?: assets
}

interface ChatDetails {
  messages: Array<ChatMessages>
  participants: Array<chatParticipants>
}

type ChatList = {
  userCode: string
  chatCode: string
  groupCode: string
  groupName: string
  senderCode: string
  senderFullName: string
  recipientCode: string
  recipientFullName: string
  messageContent: string
  messageSeen: boolean
  messageCreatedAt: Date
  assetName: string
  assetMimeType: string
}
