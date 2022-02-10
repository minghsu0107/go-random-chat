var userIDStoreKey = "rt:userid"
var userNameStoreKey = "rt:username"
var channelIDStoreKey = "rt:channelid"
const LEFT = "left"
const RIGHT = "right"

const EVENT_TEXT = 0
const EVENT_ACTION = 1

var USER_ID = ""
if (localStorage.getItem(userIDStoreKey) !== null) {
    USER_ID = localStorage.getItem(userIDStoreKey)
} else {
    window.location.href = '/'
}
var CHANNEL_ID = ""
if (localStorage.getItem(channelIDStoreKey) !== null) {
    CHANNEL_ID = localStorage.getItem(channelIDStoreKey)
} else {
    window.location.href = '/'
}
var USER_IMG = getUserImageURL(USER_ID)

var ID2NAME = {}
var USER_NAME = ""
if (localStorage.getItem(userNameStoreKey) !== null) {
    USER_NAME = localStorage.getItem(userNameStoreKey)
    ID2NAME[USER_ID] = USER_NAME
}
var ONLINE_USERS = new Set()

var chatUrl = "ws://" + window.location.host + "/api/chat?uid=" + USER_ID + "&cid=" + CHANNEL_ID
ws = new WebSocket(chatUrl)

var chatroom = document.getElementsByClassName("msger-chat")
var text = document.getElementById("msg")
var send = document.getElementById("send")
var leave = document.getElementById("leave")


var timeout = setTimeout(function () { }, 0)
var userTypingID = 'usertyping'
var peerTypingID = 'peertyping'
var isTyping = false
text.addEventListener('keypress', function () {
    clearTimeout(timeout)
    if (!isTyping) {
        insertMsg(getTypingMessage(USER_ID, RIGHT, userTypingID), chatroom[0])
        sendActionMessage("istyping")
    }
    isTyping = true
    timeout = setTimeout(function () {
        let el = document.getElementById(userTypingID)
        if (el !== null) {
            el.remove()
        }
        sendActionMessage("endtyping")
        isTyping = false
    }, 500)
})
send.onclick = function (e) {
    sendTextMessage()
}
leave.onclick = async function (e) {
    var result = confirm("Are you sure you want to leave?")
    if (result) {
        try {
            await deleteChannel()
            localStorage.removeItem(channelIDStoreKey)
            window.location.reload()
        } catch (err) {
            console.log(`Error: ${err}`)
        }
    }
}
text.onkeydown = function (e) {
    if (text.value === "\n") {
        text.value = ""
    }
    if (e.keyCode === 13 && !e.shiftKey) {
        sendTextMessage()
    }
}
ws.addEventListener('open', async function (e) {
    if (USER_NAME === "") {
        try {
            await getUserName()
        } catch (err) {
            console.log(`Error: ${err}`)
        }
    }
    try {
        await getAllChannelUserNames()
        await fetchMessages()
    } catch (err) {
        console.log(`Error: ${err}`)
    }
    document.getElementById("msg").disabled = false
})
ws.addEventListener('message', async function (e) {
    var m = JSON.parse(e.data)
    if (m.event === EVENT_ACTION) {
        switch (m.payload) {
            case "waiting":
            case "joined":
            case "offline":
                try {
                    await updateOnlineUsers()
                } catch (err) {
                    console.log(`Error: ${err}`)
                }
                break
            case "endtyping":
                let el = document.getElementById(peerTypingID)
                if (el !== null) {
                    el.remove()
                }
                break
            case "leaved":
                try {
                    await updateOnlineUsers()
                } catch (err) {
                    console.log(`Error: ${err}`)
                }
                localStorage.removeItem(channelIDStoreKey)
                ws.close()
                break
        }
    }
    var msg = await processMessage(m)
    if (msg !== "") {
        if (m.event === EVENT_TEXT) {
            let el = (m.user_id === USER_ID) ? document.getElementById(userTypingID) : document.getElementById(peerTypingID)
            if (el !== null) {
                el.remove()
            }
            sendBrowserNotification("You got a new message")
        }
        insertMsg(msg, chatroom[0])
        if (m.event === EVENT_ACTION && m.payload === "leaved") {
            insertMsg(getReturnHomeMessage(), chatroom[0])
        }
    }
})
ws.addEventListener('close', function (e) {
    document.getElementById("headstatus").innerHTML = `
    <div id="headstatus"><i class="fas fa-circle icon-red"></i>&nbsp;disconnected</div>
    `
    document.getElementById("msg").disabled = true
})
window.onbeforeunload = function () {
    ws.onclose = function () { }; // disable onclose handler first
    ws.close();
};

window.addEventListener('load', function () {
    Notification.requestPermission(function (status) {
        // This allows to use Notification.permission with Chrome/Safari
        if (Notification.permission !== status) {
            Notification.permission = status
        }
    })
})

function sendBrowserNotification(msg) {
    if (Notification && Notification.permission === "granted") {
        var n = new Notification(msg)
    }
    else if (Notification && Notification.permission !== "denied") {
        Notification.requestPermission(function (status) {
            if (Notification.permission !== status) {
                Notification.permission = status
            }
            // If the user said okay
            if (status === "granted") {
                var n = new Notification(msg)
            }
        })
    }
}

async function getAllChannelUserNames() {
    return fetch(`/api/users?cid=${CHANNEL_ID}`)
        .then((response) => {
            return response.json()
        })
        .then(async (result) => {
            for (const userID of result.user_ids) {
                if ((userID !== USER_ID) && !(userID in ID2NAME)) {
                    await setPeerName(userID)
                }
            }
        })
}

async function fetchMessages() {
    let response = await fetch(`/api/channel/${CHANNEL_ID}/messages`)
    let result = await response.json()
    for (const message of result.messages) {
        var msg = await processMessage(message)
        insertMsg(msg, chatroom[0])
    }
    if (result.messages.length === 0) {
        insertMsg(getActionMessage("Matched!"), chatroom[0])
    }
}

async function processMessage(m) {
    if (!(m.user_id in ID2NAME)) {
        await setPeerName(m.user_id)
    }
    var msg = ""
    switch (m.event) {
        case EVENT_TEXT:
            if (m.user_id === USER_ID) {
                msg = getTextMessage(USER_ID, RIGHT, m.payload, m.time)
            } else {
                msg = getTextMessage(m.user_id, LEFT, m.payload, m.time)
            }
            break
        case EVENT_ACTION:
            var actionMsg = ""
            switch (m.payload) {
                case "waiting":
                    break
                case "joined":
                    if (m.user_id !== USER_ID) {
                        actionMsg = ID2NAME[m.user_id] + " joined"
                    }
                    break
                case "offline":
                    break
                case "leaved":
                    if (m.user_id !== USER_ID) {
                        actionMsg = ID2NAME[m.user_id] + " leaved, channel closed"
                    }
                    break
                case "istyping":
                    if (m.user_id !== USER_ID) {
                        msg = getTypingMessage(m.user_id, LEFT, peerTypingID)
                    }
                    break
            }
            if (actionMsg !== "") {
                msg = getActionMessage(actionMsg)
            }
            break
    }
    return msg
}

async function getUserName() {
    return fetch(`/api/user/${USER_ID}/name`)
        .then((response) => {
            return response.json()
        })
        .then((result) => {
            USER_NAME = result.name
            ID2NAME[USER_ID] = USER_NAME
        })
}
async function setPeerName(peerID) {
    return fetch(`/api/user/${peerID}/name`)
        .then((response) => {
            return response.json()
        })
        .then((result) => {
            ID2NAME[peerID] = result.name
        })
}

async function updateOnlineUsers() {
    return fetch(`/api/users/online?cid=${CHANNEL_ID}`)
        .then((response) => {
            return response.json()
        })
        .then(async (result) => {
            var curOnlineUsers = new Set()
            for (const userID of result.user_ids) {
                var name = ""
                if (userID in ID2NAME) {
                    name = ID2NAME[userID]
                } else {
                    await setPeerName(userID)
                }
                name = ID2NAME[userID]
                curOnlineUsers.add(JSON.stringify(
                    {
                        id: userID,
                        name: name
                    }
                ))
            }
            ONLINE_USERS = curOnlineUsers
            var onlineMsg = ""
            var youMsg = ""
            for (var onlinerUserStr of ONLINE_USERS) {
                var onlineUser = JSON.parse(onlinerUserStr)
                if (onlineUser.id === USER_ID) {
                    youMsg = ", you"
                    continue
                }
                if (onlineMsg !== "") {
                    onlineMsg += ", "
                }
                onlineMsg += onlineUser.name
            }
            if (youMsg !== "") {
                if (onlineMsg === "") {
                    onlineMsg = "only you"
                } else {
                    onlineMsg += youMsg
                }
            }
            document.getElementById("headstatus").innerHTML = `
                <div id="headstatus" style="font-size: 1rem;"><i class="fas fa-circle icon-green"></i>&nbsp;${onlineMsg}</div>
                `
        })
}

async function deleteChannel() {
    return fetch(`/api/channel/${CHANNEL_ID}?delby=${USER_ID}`, {
        method: 'DELETE'
    })
}

function getUserImageURL(userID) {
    return "https://avatars.dicebear.com/api/pixel-art/" + userID + ".svg"
}

function onlySpaces(str) {
    return str.trim().length === 0
}

function sendTextMessage() {
    if (!onlySpaces(text.value)) {
        const d = new Date()
        var time = `${d.getFullYear()}/${d.getMonth()}/${d.getDay()} ${String(d.getHours()).padStart(2, "0")}:${String(d.getMinutes()).padStart(2, "0")}`
        ws.send(JSON.stringify({
            "event": EVENT_TEXT,
            "channel_id": CHANNEL_ID,
            "user_id": USER_ID,
            "payload": text.value,
            "time": time
        }))
        text.value = ""
    }
}

function sendActionMessage(action) {
    ws.send(JSON.stringify({
        "event": EVENT_ACTION,
        "channel_id": CHANNEL_ID,
        "user_id": USER_ID,
        "payload": action,
    }))
}

function getActionMessage(msg) {
    var msg = `<br><div class="msg-left">${msg}</div><br>`
    return msg
}

function getTextMessage(userID, side, text, time) {
    var msg = `
    <div class="msg ${side}-msg">
      <div class="msg-img" style="background-image: url(${getUserImageURL(userID)})"></div>

      <div class="msg-bubble" style="min-width: 125px">
        <div class="msg-info">
          <div class="msg-info-name">${ID2NAME[userID]}</div>
          <div class="msg-info-time">${time.split(' ')[1]}</div>
        </div>

        <div class="msg-text" style="max-width: 15em;overflow-wrap: break-word;">${urlify(text).replace(/(?:\r|\n|\r\n)/g, '<br>')}</div>
      </div>
    </div>
    `
    return msg
}

function getTypingMessage(userID, side, id) {
    return `
    <div class="msg ${side}-msg" id="${id}">
        <div class="msg-img" style="background-image: url(${getUserImageURL(userID)})"></div>
        <div class="chat-bubble">
            <div class="typing">
                <div class="dot"></div>
                <div class="dot"></div>
                <div class="dot"></div>
            </div>
        </div>
    </div>
    `
}

function getReturnHomeMessage() {
    return `
    <div class="msg-left"><a href="/" style="text-decoration: none; color: #1a75ff">Back</a></div><br>
    `
}

function insertMsg(msg, domObj) {
    domObj.insertAdjacentHTML("beforeend", msg)
    domObj.scrollTop += 500
}

function urlify(text) {
    var urlRegex = /(https?:\/\/[^\s]+)/g;
    return text.replace(urlRegex, function (url) {
        return '<a href="' + url + '">' + url + '</a>';
    })
    // or alternatively
    // return text.replace(urlRegex, '<a href="$1">$1</a>')
}