// Doorlock Handle

const mode = {
    idle: 0,
    checking: 1,
    authentication: 2
}

let options = new MQTTConfig()
options.lastWill = false
options.lastWillMessage = ""
options.lastWillRetain = false
options.Hostname = "192.168.100.10"
options.Port = 1883

let client = new MQTTWrapper(options)
client.Subscribe("doorlock/+/status", statusHandle)
client.Subscribe("doorlock/+/read", readHandle)

let doorlock = new Map()

function statusHandle(topic, payload) {
    let tokens = topic.split('/')
    let id = tokens[1]
    let status = payload !== "false";
    doorlock.set(id, {status: status, mode: mode.checking})
}

function readHandle(topic, payload) {
    let tokens = topic.split('/')
    let id = tokens[1]
    lock = doorlock.get(id)
    readObj = JSON.parse(payload)

    if (lock.mode === mode.checking) {
        check(id, readObj)
    } else if (lock.mode === mode.authentication) {
        authenticate(id, readObj)
    }
}

function check(id, readObj) {
    try {
        let obj = JSON.parse(ReadFile(readObj.uid + ".json"))
        if (obj.uid === readObj.uid && obj.data === readObj.data)
            openDoor(id)
    } catch (e) {

    }
}

function authenticate(id, readObj) {
    try {
        obj = {uid: readObj.uid, data: RandomBase64Bytes(48)}
        WriteFile(readObj.uid, JSON.stringify(obj))
        client.Publish("doorlock/" + id + "/write/data", obj.data, false)
    } catch (e) {

    }
}

function openDoor(id) {
    client.Publish("doorlock/" + id + "/open", "true", false)
}

function close() {
    doorlock.forEach((key, value) => {
        client.Unsubscribe("doorlock/" + key + "/status")
        client.Unsubscribe("doorlock/" + key + "/read")
    })
    doorlock.clear()
}