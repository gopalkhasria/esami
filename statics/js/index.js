var token;
function inizializzo() {
    token = document.cookie.split('=');
    token = token[1];
    const socket = new WebSocket("ws://localhost:5000/ws");
    socket.onopen = () => {
        console.log("Successfully Connected");
        socket.send("Hi From the Client!")
    };
    socket.onmessage = (msg) => {
        console.log(msg.data);
    }
    socket.onerror = error => {
        console.log("Socket Error: ", error);
    };
}
