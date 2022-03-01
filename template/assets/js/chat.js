window.mobileCheck = function () {
    let check = false;
    (function (a) { if (/(android|bb\d+|meego).+mobile|avantgo|bada\/|blackberry|blazer|compal|elaine|fennec|hiptop|iemobile|ip(hone|od)|iris|kindle|lge |maemo|midp|mmp|mobile.+firefox|netfront|opera m(ob|in)i|palm( os)?|phone|p(ixi|re)\/|plucker|pocket|psp|series(4|6)0|symbian|treo|up\.(browser|link)|vodafone|wap|windows ce|xda|xiino/i.test(a) || /1207|6310|6590|3gso|4thp|50[1-6]i|770s|802s|a wa|abac|ac(er|oo|s\-)|ai(ko|rn)|al(av|ca|co)|amoi|an(ex|ny|yw)|aptu|ar(ch|go)|as(te|us)|attw|au(di|\-m|r |s )|avan|be(ck|ll|nq)|bi(lb|rd)|bl(ac|az)|br(e|v)w|bumb|bw\-(n|u)|c55\/|capi|ccwa|cdm\-|cell|chtm|cldc|cmd\-|co(mp|nd)|craw|da(it|ll|ng)|dbte|dc\-s|devi|dica|dmob|do(c|p)o|ds(12|\-d)|el(49|ai)|em(l2|ul)|er(ic|k0)|esl8|ez([4-7]0|os|wa|ze)|fetc|fly(\-|_)|g1 u|g560|gene|gf\-5|g\-mo|go(\.w|od)|gr(ad|un)|haie|hcit|hd\-(m|p|t)|hei\-|hi(pt|ta)|hp( i|ip)|hs\-c|ht(c(\-| |_|a|g|p|s|t)|tp)|hu(aw|tc)|i\-(20|go|ma)|i230|iac( |\-|\/)|ibro|idea|ig01|ikom|im1k|inno|ipaq|iris|ja(t|v)a|jbro|jemu|jigs|kddi|keji|kgt( |\/)|klon|kpt |kwc\-|kyo(c|k)|le(no|xi)|lg( g|\/(k|l|u)|50|54|\-[a-w])|libw|lynx|m1\-w|m3ga|m50\/|ma(te|ui|xo)|mc(01|21|ca)|m\-cr|me(rc|ri)|mi(o8|oa|ts)|mmef|mo(01|02|bi|de|do|t(\-| |o|v)|zz)|mt(50|p1|v )|mwbp|mywa|n10[0-2]|n20[2-3]|n30(0|2)|n50(0|2|5)|n7(0(0|1)|10)|ne((c|m)\-|on|tf|wf|wg|wt)|nok(6|i)|nzph|o2im|op(ti|wv)|oran|owg1|p800|pan(a|d|t)|pdxg|pg(13|\-([1-8]|c))|phil|pire|pl(ay|uc)|pn\-2|po(ck|rt|se)|prox|psio|pt\-g|qa\-a|qc(07|12|21|32|60|\-[2-7]|i\-)|qtek|r380|r600|raks|rim9|ro(ve|zo)|s55\/|sa(ge|ma|mm|ms|ny|va)|sc(01|h\-|oo|p\-)|sdk\/|se(c(\-|0|1)|47|mc|nd|ri)|sgh\-|shar|sie(\-|m)|sk\-0|sl(45|id)|sm(al|ar|b3|it|t5)|so(ft|ny)|sp(01|h\-|v\-|v )|sy(01|mb)|t2(18|50)|t6(00|10|18)|ta(gt|lk)|tcl\-|tdg\-|tel(i|m)|tim\-|t\-mo|to(pl|sh)|ts(70|m\-|m3|m5)|tx\-9|up(\.b|g1|si)|utst|v400|v750|veri|vi(rg|te)|vk(40|5[0-3]|\-v)|vm40|voda|vulc|vx(52|53|60|61|70|80|81|83|85|98)|w3c(\-| )|webc|whit|wi(g |nc|nw)|wmlb|wonu|x700|yas\-|your|zeto|zte\-/i.test(a.substr(0, 4))) check = true; })(navigator.userAgent || navigator.vendor || window.opera);
    return check;
};
var userIDStoreKey = "rc:userid"
var userNameStoreKey = "rc:username"
var accessTokenKey = "rc:accesstoken"
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
var ACCESS_TOKEN = ""
if (localStorage.getItem(accessTokenKey) !== null) {
    ACCESS_TOKEN = localStorage.getItem(accessTokenKey)
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

var chatUrl = "ws://" + window.location.host + "/api/chat?uid=" + USER_ID + "&access_token=" + ACCESS_TOKEN
ws = new WebSocket(chatUrl)

var chatroom = document.getElementsByClassName("msger-chat")
var text = document.getElementById("msg")
var send = document.getElementById("send")
var leave = document.getElementById("leave")


var timeout = setTimeout(function () { }, 0)
var userTypingID = 'usertyping'
var peerTypingID = 'peertyping'
var isTyping = false
text.addEventListener('keyup', function () {
    clearTimeout(timeout)
    if (!isTyping) {
        insertMsg(getTypingMessage(USER_ID, RIGHT, userTypingID), chatroom[0], true)
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

send.addEventListener("pointerdown", function (e) {
    send.style.color = "#0a1869"
})
send.addEventListener("pointerup", function (e) {
    sendTextMessage()
    text.setAttribute("rows", 1)
    send.style.color = "#25A3FF"
})
leave.onclick = async function (e) {
    var result = confirm("Are you sure you want to leave?")
    if (result) {
        try {
            await deleteChannel()
            localStorage.removeItem(accessTokenKey)
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
    if (window.mobileCheck()) {
        if (e.keyCode === 13) {
            text.setAttribute("rows", parseInt(text.getAttribute("rows"), 10) + 1)
        }
    } else if (e.keyCode === 13) {
        if (!e.shiftKey) {
            sendTextMessage()
            text.setAttribute("rows", 1)
        } else {
            text.setAttribute("rows", parseInt(text.getAttribute("rows"), 10) + 1)
        }
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
                localStorage.removeItem(accessTokenKey)
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
            if (!window.mobileCheck) {
                sendBrowserNotification("You got a new message")
            }
        }
        var isSelf = (m.user_id === USER_ID)
        insertMsg(msg, chatroom[0], isSelf)
        if (m.event === EVENT_ACTION && m.payload === "leaved") {
            insertMsg(getReturnHomeMessage(), chatroom[0], isSelf)
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
    return fetch(`/api/users`, {
        method: 'GET',
        headers: new Headers({
            'Authorization': 'Bearer ' + ACCESS_TOKEN
        })
    })
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
    let response = await fetch(`/api/channel/messages`, {
        method: 'GET',
        headers: new Headers({
            'Authorization': 'Bearer ' + ACCESS_TOKEN
        })
    })
    let result = await response.json()
    for (const message of result.messages) {
        var msg = await processMessage(message)
        insertMsg(msg, chatroom[0], true)
    }
    if (result.messages.length === 0) {
        insertMsg(getActionMessage("Matched!"), chatroom[0], true)
    }
}

async function processMessage(m) {
    if (!(m.user_id in ID2NAME)) {
        await setPeerName(m.user_id)
    }
    var msg = ""
    switch (m.event) {
        case EVENT_TEXT:
            const d = new Date(m.time)
            var time = `${d.getFullYear()}/${d.getMonth() + 1}/${d.getDate()} ${String(d.getHours()).padStart(2, "0")}:${String(d.getMinutes()).padStart(2, "0")}`
            if (m.user_id === USER_ID) {
                msg = getTextMessage(USER_ID, RIGHT, m.payload, time)
            } else {
                msg = getTextMessage(m.user_id, LEFT, m.payload, time)
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
    return fetch(`/api/users/online`, {
        method: 'GET',
        headers: new Headers({
            'Authorization': 'Bearer ' + ACCESS_TOKEN
        })
    })
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
    return fetch(`/api/channel?delby=${USER_ID}`, {
        method: 'DELETE',
        headers: new Headers({
            'Authorization': 'Bearer ' + ACCESS_TOKEN
        })
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
        ws.send(JSON.stringify({
            "event": EVENT_TEXT,
            "user_id": USER_ID,
            "payload": text.value,
        }))
        text.value = ""
    }
}

function sendActionMessage(action) {
    ws.send(JSON.stringify({
        "event": EVENT_ACTION,
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

function insertMsg(msg, domObj, isSelf) {
    domObj.insertAdjacentHTML("beforeend", msg)
    if (isSelf) {
        domObj.scrollTop = domObj.scrollHeight
    } else {
        if (domObj.scrollHeight - 1.2 * domObj.offsetHeight <= domObj.scrollTop) {
            domObj.scrollTop = domObj.scrollHeight
        }
    }
    if (text.value === "\n") {
        text.value = ""
    }
}

function urlify(text) {
    var urlRegex = /(https?:\/\/[^\s]+)/g;
    return text.replace(urlRegex, function (url) {
        return '<a href="' + url + '">' + url + '</a>';
    })
    // or alternatively
    // return text.replace(urlRegex, '<a href="$1">$1</a>')
}

function auto_grow(element) {
    element.style.height = "5px";
    element.style.height = (element.scrollHeight) + "px";
}
