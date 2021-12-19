import '../css/Button.css'
import '../css/ComponentPage.css'
import '../css/shadow.css'
import React from 'react'
import getComponentForType from "./Util";

class Data extends React.Component {
    timeout = 250; // Initial timeout duration as a class variable

    constructor(props) {
        super(props);
        this.state = {
            name: props.name,
            ws: null,
            child: null,
            updateRequest: props.updateRequest,
            updateSocket: props.updateSocket
        }
    }

    componentDidMount() {
        this.connect()

        fetch("http://" + window.location.hostname + ":1337/" + this.state.updateRequest, {
            credentials: "include"
        }).then(data => data.json()).then(data => {
            this.setState({child: data})
        })
    }

    /**
     * @function connect
     * This function establishes the connect with the websocket and also ensures constant reconnection if connection closes
     */
    connect = () => {
        let ws = new WebSocket("ws://" + window.location.hostname + ":1337/" + this.state.updateSocket);
        let that = this; // cache the this
        let connectInterval;

        // websocket onopen event listener
        ws.onopen = () => {
            this.setState({ws: ws});

            that.timeout = 250; // reset timer to 250 on open of websocket connection
            clearTimeout(connectInterval); // clear Interval on on open of websocket connection
        };

        // websocket onclose event listener
        ws.onclose = e => {
            console.log(e.reason)
            that.timeout = that.timeout + that.timeout; //increment retry interval
            connectInterval = setTimeout(this.check, Math.min(10000, that.timeout)); //call check function after timeout
        };

        // websocket onerror event listener
        ws.onerror = (error) => {
            console.error(error)
            ws.close();
        };

        ws.onmessage = (msg) => {
            this.setState({child: null})
            this.setState({child: JSON.parse(msg.data)})
        }
    };

    check = () => {
        const {ws} = this.state;
        if (!ws || ws.readyState === WebSocket.CLOSED) this.connect(); //check if websocket instance is closed, if so call `connect` function.
    };

    render() {
        if (this.state.child === null) {
            return null;
        }
        return getComponentForType(this.state.child);
    }

}

export default Data