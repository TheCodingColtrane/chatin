const getAllCookies = () => {
    const rawCookies = document.cookie
    const splittedCookies = rawCookies.split('=')
    const cookies = [{name: "", value: ""}]
    if(splittedCookies.length === 1) return cookies

    for(let i =0; i < splittedCookies.length; i++){        
        if(i + 2 <= splittedCookies.length) {
            if(i === 0){
            cookies[i].name  = splittedCookies[i]
            cookies[i].value  = splittedCookies[i + 1].replace(";", "")}
            i = i + 1
            continue;
            }
            cookies.push({name: splittedCookies[i], value: splittedCookies[i + 1].replace(";", "")})      
            i = i + 1
            continue;
        }

        return cookies
}

const findCookie = (name: string) => {
    return getAllCookies().find(c => c.name === name)
}

const createCookie = (name: string, value: string, exp: number) => {
    document.cookie = `${name}=${value};Path=/; Max-Age=${exp}`
}

const deleteCookie = (name: string, value: string) => {
    document.cookie = `${name}=${value};Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT}`

}


export {
    createCookie,
    deleteCookie,
    findCookie,
    getAllCookies
}


