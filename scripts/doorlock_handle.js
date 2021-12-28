// Doorlock Handle

const mode = {
    idle: 0,
    checking: 1,
    authentication: 2
}
const doorlocks = new Map()

function newEspObj(id) {
    return {
        id: id,
        status: false,
        mode: mode.checking,
        container: null,
        data: null
    }
}

function setStatus(espObj, status) {
    espObj.status = status

    if (status === true) {
        if (espObj.container === null) {
            createContainer(espObj)
        }
        AddContainer(espObj.container)
    } else if (status === false) {
        if (espObj.container !== null) {
            RemoveContainer(espObj.container)
        }
    }
}

function createContainer(espObj) {
    espObj.container = new Container("doorlock" + espObj.id, "Lock " + espObj.id, (user) => {
    })

    let btnCheck = new Button("doorlock" + espObj.id + "check", "Check", (user) => {
        espObj.mode = mode.checking
        data.Update(new Alert("doorlock" + espObj.id + "alert", "Changed to CHECK", "success"))
        ExecuteAfterMs(20000, () => {
            data.Update(new Alert("doorlock" + espObj.id + "alert", "", "success"))
        })
    })
    let btnAuth = new Button("doorlock" + espObj.id + "auth", "Authenticate", (user) => {
        espObj.mode = mode.authentication
        data.Update(new Alert("doorlock" + espObj.id + "alert", "Changed to AUTH", "success"))
        ExecuteAfterMs(20000, () => {
            data.Update(new Alert("doorlock" + espObj.id + "alert", "", "success"))
        })
    })
    let btnIdle = new Button("doorlock" + espObj.id + "idle", "Idle", (user) => {
        espObj.mode = mode.idle
        data.Update(new Alert("doorlock" + espObj.id + "alert", "Changed to IDLE", "info"))
        ExecuteAfterMs(20000, () => {
            data.Update(new Alert("doorlock" + espObj.id + "alert", "", "success"))
        })
    })

    let check1 = new Checkbox("check1", "Check Me")
    check1.SetChangeCallback((username, state) => {
        Print(username + ": " + state)
    })

    let textfield = new TextField("text1", "Text", (user, text) => {
        Print(user + ": " + text)
    })

    let data = new Data("doorlock" + espObj.id + "data", new Alert("doorlock" + espObj.id + "alert", "", "info"))
    espObj.data = data

    espObj.container.Add(btnCheck)
    espObj.container.Add(btnAuth)
    espObj.container.Add(btnIdle)
    espObj.container.Add(textfield)
    espObj.container.Add(check1)
    espObj.container.Add(data)
}

function check(espObj, readObj) {
    try {
        let obj = JSON.parse(ReadFile(readObj.uid + ".json"))
        if (obj.uid === readObj.uid && obj.data === readObj.data)
            openDoor(espObj)
    } catch (e) {

    }
}

function authenticate(espObj, readObj) {
    try {
        let obj = {uid: readObj.uid, data: RandomBase64Bytes(48)}
        WriteFile(readObj.uid + ".json", JSON.stringify(obj))
        client.Publish("doorlock/" + espObj.id + "/write/data", obj.data, false)
    } catch (e) {
    }
}

function openDoor(espObj) {
    client.Publish("doorlock/" + espObj.id + "/open", "true", false)
}

function addIfNotExists(id, status) {
    if (!doorlocks.has(id)) {
        let handle = newEspObj(id)
        setStatus(handle, status)
        doorlocks.set(id, handle)
    }
}

function statusHandle(topic, payload) {
    let tokens = topic.split('/')
    let id = tokens[1]
    let status = payload !== "false";
    addIfNotExists(id, false)
    setStatus(doorlocks.get(id), status)
}

function readHandle(topic, payload) {
    let tokens = topic.split('/')
    let id = tokens[1]
    let lock = doorlocks.get(id)
    let readObj = JSON.parse(payload)

    if (lock.mode === mode.checking) {
        check(lock, readObj)
    } else if (lock.mode === mode.authentication) {
        authenticate(lock, readObj)
    }
}


function close() {
    client.Unsubscribe("doorlock/+/status")
    client.Unsubscribe("doorlock/+/read")
    doorlocks.clear()
}

let options = new MQTTConfig()
options.lastWill = true
options.lastWillMessage = "false"
options.lastWillRetain = true
options.lastWillTopic = "doorlock/status"
options.Hostname = "192.168.100.10"
options.Port = 1883

let client = new MQTTWrapper(options)
client.Subscribe("doorlock/+/status", statusHandle)
client.Subscribe("doorlock/+/read", readHandle)

