import endpoints from ".."

const baseUrl = "http://localhost:3000"
const fileServer = "http://localhost:8001"
const users = baseUrl + "/users"
const messages = baseUrl + "/chats"
const contacts = baseUrl + "/contacts"
const auth = baseUrl + "/auth"
let uri = ""

const backend = {
    users: {
        get: async function<T>(params = ""): Promise<T | undefined> {
            try {
                params ? uri = users + params : users
                const data = await endpoints.get<T>(uri)
                return data
            } catch (error) {
                console.log(error as Error)
            }
        },
        post: async function<T>(data: T, params = ""): Promise<T | undefined> {
            try {
                params ? uri = users + params : users
                const result = await endpoints.post<T>(uri, data)
                return result
            } catch (error) {
                console.log(error as Error)
            }
        } 
    },

    messages: {
        get: async function<T>(params = ""): Promise<T | undefined> {
            try {
                params ? uri = messages + params : messages
                const data = await endpoints.get<T>(uri)
                return data
            } catch (error) {
                console.log(error as Error)
            }
        },
        post: async function<T>(data: T, params = ""): Promise<T | undefined> {
            try {
                params ? uri = messages + params : messages
                const result = await endpoints.post<T>(uri, data)
                return result
            } catch (error) {
                console.log(error as Error)
            }
        } 
    },

    contacts: {
        get: async function<T>(params = ""): Promise<T | undefined> {
            try {
                params ? uri = contacts + params : contacts
                const data = await endpoints.get<T>(uri)
                return data
            } catch (error) {
                console.log(error as Error)
            }
        },
        post: async function<T>(data: T, params = ""): Promise<T | undefined> {
            try {
                params ? uri = contacts + params : contacts
                const result = await endpoints.post<T>(uri, data)
                return result
            } catch (error) {
                console.log(error as Error)
            }
        } 
    },

    auth: {
        post: async function<T>(data: T, params = ""): Promise<T | undefined> {
            try {
                uri = params ?  auth + params : auth
                const result = await endpoints.post<T>(uri, data)
                return result
            } catch (error) {
                console.log(error as Error)
            }
        } 
    }
    
}
 




export {
    backend,
    fileServer
}