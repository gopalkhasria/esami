var token;
var amount = 0;
var adress;
var myTransaction = [];
var myOutputs = [];
var socket;
startWebsocket();
function startWebsocket() {
    socket = new WebSocket("ws://localhost:5000/ws");
    socket.onopen = () => {
        console.log("Successfully Connected");
        socket.send("Hi From the Client!")
    };
    socket.onmessage = (msg) => {
        start(msg)
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
function inizializzo() {
    token = document.cookie.split('=');
    token = token[1];
}

function start(msg) {
    adress = document.getElementById("adress").innerText;
    data = JSON.parse(msg.data);
    if (data.azione == 1) {
        amount = 0;
        var html = '';
        var html2 = ''
        for (var i = 0; i < data.transaction.length; i++) {
            html += '<div class="card"><a href="/transaction?id=' + data.transaction[i].id + '">' + data.transaction[i].hash + '</a></diV>';
            if (data.transaction[i].output.pkScript == adress) {
                html2 += '<div class="card"><a href="/transaction?id=' + data.transaction[i].id + '">' + data.transaction[i].hash + '</a></diV>';
                myTransaction.push(data.transaction[i].hash);
                myOutputs.push(data.transaction[i].output);
                if (data.transaction[i].output) {
                    amount += parseInt(data.transaction[i].output.amount);
                }
            }
        }
        document.getElementById("amount").innerText = amount;
        document.getElementById("transactions").innerHTML = html;
        document.getElementById("MyTransactions").innerHTML = html2;
    }
    if (socket.readyState !== socket.OPEN) {
        startWebsocket();
    }
    console.log(myOutputs);
}

function swithscreen(screen) {
    switch (screen) {
        case 1:
            document.getElementById("transactions").style.display = 'block';
            document.getElementById("MyTransactions").style.display = 'none';
            document.getElementById("sendBit").style.display = 'none';
            break;
        case 2:
            document.getElementById("transactions").style.display = 'none';
            document.getElementById("MyTransactions").style.display = 'block';
            document.getElementById("sendBit").style.display = 'none';
            break;
        case 3:
            document.getElementById("transactions").style.display = 'none';
            document.getElementById("MyTransactions").style.display = 'none';
            document.getElementById("sendBit").style.display = 'block';
            break;
    }
}