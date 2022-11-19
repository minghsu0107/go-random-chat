var accessTokenKey = "rc:accesstoken"

var questions = [
    { question: "What's your name?", pattern: /^.{1,15}$/ },
    // { question: "How old are you?", pattern: /^100|[1-9]?\d$/ }
]
var tTime = 100  // transition transform time from #register in ms
var wTime = 200  // transition width time from #register in ms
var eTime = 1000 // transition width time from inputLabel in ms
var position = 0
var ws

var isLogin = false
async function getUserInfo() {
    return fetch(`/api/user/me`, {
        method: 'GET'
    })
        .then((response) => {
            if (response.status !== 200) {
                throw Error(response.statusText)
            }
            isLogin = true
        })
        .catch((error) => {
            console.log(`Error: ${error}`)
            isLogin = false
        })
}
async function check() {
    await getUserInfo()
    if (!isLogin) {
        (function () {
            localStorage.removeItem(accessTokenKey)
            putQuestion()
    
            progressButton.addEventListener('click', validate)
            inputField.addEventListener('keyup', function (e) {
                transform(0, 0) // ie hack to redraw
                if (e.keyCode == 13) validate()
            })
        }())
    } else if (localStorage.getItem(accessTokenKey) === null) {
        inputContainer.style.opacity = 0
        inputProgress.style.transition = 'none'
        inputProgress.style.width = 0
        done()
    } else {
        window.location.href = '/chat'
    }
} 
check()

// load the next question
function putQuestion() {
    inputLabel.innerHTML = questions[position].question
    inputField.value = ''
    inputField.type = questions[position].type || 'text'
    inputField.focus()
    showCurrent()
}

// when all the questions have been answered
function done() {
    // remove the box if there is no next question
    register.className = 'close'

    // add the h1 at the end with the welcome text
    var h1 = document.createElement('h1')
    var button = document.createElement("button")
    button.setAttribute("id", "startbutton");
    button.className = "start-btn"
    button.innerHTML = 'Start'
    button.onclick = function () {
        button.disabled = true
        button.innerHTML = `
        <p class="saving">Matching<span>.</span><span>.</span><span>.</span></p>
        `
        button.style.cursor = 'default'
        match()
    }
    h1.appendChild(button)
    setTimeout(function () {
        register.parentElement.appendChild(h1)
        setTimeout(function () { h1.style.opacity = 1 }, 50)
    }, eTime)
}

async function createUser(username) {
    var data = {
        name: username,
    }
    return fetch(`/api/user`, {
        body: JSON.stringify(data),
        method: 'POST'
    })
        .then((response) => {
            if (response.status !== 201) {
                throw Error(response.statusText)
            }
        })
        .catch((error) => {
            console.log(`Error: ${error}`)
        })
}
function match() {
    var protocol
    var loc = window.location
    if (loc.protocol === "https:") {
        protocol = "wss:"
    } else {
        protocol = "ws:"
    }
    var matchUrl = protocol + "//" + window.location.host + "/api/match"
    ws = new WebSocket(matchUrl)
    ws.addEventListener('message', function (e) {
        var result = JSON.parse(e.data)
        if (result.channel_id !== "" && result.access_token !== "") {
            localStorage.setItem(accessTokenKey, result.access_token)
            ws.close()
            window.location.href = '/chat'
        }
    })
}
window.onbeforeunload = function () {
    ws.onclose = function () { }; // disable onclose handler first
    ws.close();
};

// when submitting the current question
async function validate() {

    // set the value of the field into the array
    questions[position].value = inputField.value

    // check if the pattern matches
    if (!inputField.value.match(questions[position].pattern || /.+/)) wrong()
    else ok(async function () {
        if (position === 0) {
            await createUser(questions[0].value)
        }
        // set the progress of the background
        progress.style.width = ++position * 100 / questions.length + 'vw'

        // if there is a new question, hide current and load next
        if (questions[position]) hideCurrent(putQuestion)
        else hideCurrent(done)

    })

}

// helper
// --------------

function hideCurrent(callback) {
    inputContainer.style.opacity = 0
    inputProgress.style.transition = 'none'
    inputProgress.style.width = 0
    setTimeout(callback, wTime)
}

function showCurrent(callback) {
    inputContainer.style.opacity = 1
    inputProgress.style.transition = ''
    inputProgress.style.width = '100%'
    setTimeout(callback, wTime)
}

function transform(x, y) {
    register.style.transform = 'translate(' + x + 'px ,  ' + y + 'px)'
}

function ok(callback) {
    register.className = ''
    setTimeout(transform, tTime * 0, 0, 10)
    setTimeout(transform, tTime * 1, 0, 0)
    setTimeout(callback, tTime * 2)
}

function wrong(callback) {
    register.className = 'wrong'
    for (var i = 0; i < 6; i++) // shaking motion
        setTimeout(transform, tTime * i, (i % 2 * 2 - 1) * 20, 0)
    setTimeout(transform, tTime * 6, 0, 0)
    setTimeout(callback, tTime * 7)
}