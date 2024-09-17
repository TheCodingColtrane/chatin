import { findCookie } from "@/utils/cookies"

class Backend {

    private headers = new Headers()

    constructor(){
        this.headers.append("Content-Type", "application/json")
        this.headers.append("Authorization", "Bearer " + findCookie("Authorization")?.value)

    }
    
    async get<T>(uri: string):Promise<T | undefined > {
        try {
            const request = await fetch(uri,  { headers: this.headers, method: "GET", cache: "force-cache"})
            return await request.json() as T
        } catch (error) {
            console.log(error)
            return 
        }
   
    } 

    async post<T>(uri: string, data: T):Promise<T | undefined > {
        try {
            const request = await fetch(uri,  { headers: this.headers, method: "POST", body: JSON.stringify(data)})
            return await request.json() as T
        } catch (error) {
            console.log(error)
            return 
        }
   
    } 

}


export default new Backend() 