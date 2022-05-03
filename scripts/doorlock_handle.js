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
    let data = new Data("doorlock" + espObj.id + "data", new Alert("doorlock" + espObj.id + "alert", "", "info"))
    espObj.data = data

    espObj.container = new Container("doorlock" + espObj.id, "Lock " + espObj.id, (user) => {
    }, (user) => {
        espObj.data.Update(new Alert("doorlock" + espObj.id + "alert", "", "success"))
    })

    let btnCheck = new Button("doorlock" + espObj.id + "check", "Check", (user) => {
        espObj.mode = mode.checking
        espObj.data.Update(new Alert("doorlock" + espObj.id + "alert", "Changed to CHECK", "success"))
    })
    let btnAuth = new Button("doorlock" + espObj.id + "auth", "Authenticate", (user) => {
        espObj.mode = mode.authentication
        espObj.data.Update(new Alert("doorlock" + espObj.id + "alert", "Changed to AUTH", "success"))
    })
    let btnIdle = new Button("doorlock" + espObj.id + "idle", "Idle", (user) => {
        espObj.mode = mode.idle
        espObj.data.Update(new Alert("doorlock" + espObj.id + "alert", "Changed to IDLE", "success"))
    })

    espObj.container.Add(btnCheck)
    espObj.container.Add(btnAuth)
    espObj.container.Add(btnIdle)
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
        // TODO send with tcp
        client.Publish("doorlock/" + espObj.id + "/write/data", obj.data, false)
    } catch (e) {
    }
}

function openDoor(espObj) {
    // TODO write to opener chip
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
    doorlocks.clear()
}


let ip_reader = "192.168.100.108";
let port_reader = 5677;
let ip_opener = "192.168.100.109";
let port_opener = 5656;



