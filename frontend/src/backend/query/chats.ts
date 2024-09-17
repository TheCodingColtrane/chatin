export type chatMessages = {
    code: string
	messageCode: string
	userCode: string
	username: string
	fullName: string
	seen?: boolean
	seenAt?: Date
	content:  string 
	createdAt: Date
	updatedAt: Date
	asset?: assets	
}

export type assets = {
	mimeType: string
	name: string,
	reason: number
}