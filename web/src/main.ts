import { createApp } from "vue";
import "@/assets/index.css";
import App from "./App.vue";

let flag = true

const setTokenHeaders = (e: CustomEvent) => {
    e.detail.headers["Authorization"] = "Bearer " + localStorage.getItem('token');

}

const loadToken = (e:CustomEvent) => {

    const token = e.detail.value as string;

    localStorage.setItem('token', token)
    document.cookie = "token="+token

}

const getMessageFromWebSocket = (_: CustomEvent) =>{
    const messageContent = document.getElementById('message-container') as HTMLDivElement;
    messageContent.scrollTo({
        top: messageContent.scrollHeight,
        behavior: 'smooth'
    });

    if (flag){
        setTimeout(()=>{
            const intersect = document.getElementById("first-intersection") as HTMLDivElement
            intersect.classList.remove("hidden")
        }, 300)

        flag = false
    }

}

document.addEventListener("jwt", (loadToken) as EventListener);

document.addEventListener("htmx:configRequest", (setTokenHeaders) as EventListener);

document.body.addEventListener('htmx:wsAfterMessage', (getMessageFromWebSocket as EventListener));

createApp(App).mount("#app");
