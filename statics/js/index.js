var token;
var amount = 0;
var tempAmount = 0;
var address;
var myTransaction = [];
var myOutputs = [];
var socket;
startWebsocket();
function startWebsocket() {
    socket = new WebSocket("wss://gopal-bitwallet.herokuapp.com/ws");
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
    address = document.getElementById("adress").innerText;
    data = JSON.parse(msg.data);
    var html = '';
    var html2 = '';
    var html3 = '';
    var j = 0;
    if (data.azione == 1) {
        myOutputs = [];
        amount = 0;
        tempAmount = 0;
        var tempHash = "";
        console.log(data)
        for (var i = data.transaction.length - 1; i >= 0; i--) {
            //console.log(data.transaction[i].hash)
            if (data.transaction[i].hash != tempHash) {
                html += '<div class="card"><div class="trans-fix"><a href="/transaction?id=' + data.transaction[i].id + '">' + data.transaction[i].hash + '</a><span class="dot ';
                if (data.transaction[i].block == "-1" || data.transaction[i].block == "-2") html += 'red'; else html += 'green';
                html += '"></span></div></div>'
                tempHash = data.transaction[i].hash;
            }
            //console.log(data.transaction[i].output);
            if (data.transaction[i].output.pkScript == address) {
                html2 += '<div class="card"><a href="/transaction?id=' + data.transaction[i].id + '">' + data.transaction[i].hash + '</a></diV>';
                myTransaction.push(data.transaction[i].hash);
                myOutputs.push(data.transaction[i].output);
                myOutputs[j].amount = parseFloat(myOutputs[j].amount);
                j++;
                if (!data.transaction[i].output.used) {
                    if (data.transaction[i].block == "-1" || data.transaction[i].block == "-2") {
                        tempAmount += parseFloat(data.transaction[i].output.amount);
                        console.log("gg")
                    }
                    else amount += parseFloat(data.transaction[i].output.amount);
                }
            }
        }
        for (var i = data.block.length - 1; i >= 0; i--) {
            html3 += '<div class="card"><a href="/block?id=' + data.block[i].id + '">' + data.block[i].hash + '</a></div>';
        }
        document.getElementById("amount").innerText = amount;
        document.getElementById("tempamount").innerText = tempAmount;
        document.getElementById("transactions").innerHTML = html;
        document.getElementById("MyTransactions").innerHTML = html2;
        document.getElementById("block").innerHTML = html3;
    }
    if (socket.readyState !== socket.OPEN) {
        startWebsocket();
    }
    //console.log(myOutputs);
}

function swithscreen(screen) {
    switch (screen) {
        case 1:
            document.getElementById("transactions").style.display = 'block';
            document.getElementById("MyTransactions").style.display = 'none';
            document.getElementById("sendBit").style.display = 'none';
            document.getElementById("block").style.display = 'none';
            break;
        case 2:
            document.getElementById("transactions").style.display = 'none';
            document.getElementById("MyTransactions").style.display = 'block';
            document.getElementById("sendBit").style.display = 'none';
            document.getElementById("block").style.display = 'none';
            break;
        case 3:
            document.getElementById("transactions").style.display = 'none';
            document.getElementById("MyTransactions").style.display = 'none';
            document.getElementById("sendBit").style.display = 'block';
            document.getElementById("block").style.display = 'none';
            break;
        case 4:
            document.getElementById("transactions").style.display = 'none';
            document.getElementById("MyTransactions").style.display = 'none';
            document.getElementById("sendBit").style.display = 'none';
            document.getElementById("block").style.display = 'block';
            break;
    }
}

function maketransaction() {
    var data = {
        amount: parseFloat(document.getElementById("cointosend").value),
        address: document.getElementById("address").value,
        pubkey: address,
        outputs: myOutputs
    }
    if (data.amount > amount) {
        alert("Non hai tutti questi soldi");
    } else {
        console.log(JSON.stringify(data));
        var myHeaders = new Headers();
        myHeaders.append("Authorization", token);
        myHeaders.append("Content-Type", "application/json");
        var requestOptions = {
            method: 'POST',
            headers: myHeaders,
            body: JSON.stringify(data),
            redirect: 'follow'
        };
        fetch("http://localhost:5000/makeTransaction", requestOptions)
            .then(response => response.text())
            .then(result => location.reload())
            .catch(error => console.log('error', error));
    }
}