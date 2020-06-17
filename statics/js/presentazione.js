var socket;
startWebsocket();
function startWebsocket() {
    socket = new WebSocket("wss://gopal-bitwallet.herokuapp.com/presSocket");
    //socket = new WebSocket("ws://localhost:5000/presSocket");
    socket.onopen = () => {
        console.log("Successfully Connected");
        socket.send("Hi From the Client!")
    };
    socket.onmessage = (msg) => {
        action(msg.data)
        //return false;
    }
    socket.onerror = error => {
        console.log("Socket Error: ", error);
    };

    socket.onclose = event => {
        console.log("Socket Closed Connection: ", event);
        socket.send("Client Closed!")
        socket = null
        startWebsocket()
    };
    setInterval(function () {
        if (socket.readyState !== socket.OPEN) {
            startWebsocket();
        }
    }, 3000);
}

function action(event) {
    //console.log(event)
    switch (event) {
        case '4':
            Reveal.left();
            break;
        case '3':
            Reveal.right();
            break;
        case '1':
            Reveal.up();
            break;
        case '2':
            Reveal.down();
            break;
    }
}