import '../css/ComponentPage.css'
import '../css/shadow.css'
import MuiTextField from "@mui/material/TextField"
import React from 'react'

class TextField extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            updateRequest: props.updateRequest,
            name: props.name,
            text: props.text,
            checked: false,
            content: ""
        }
    }

    getID() {
        return this.state.onUpdate + this.state.name
    }

    render() {
        return (
            <div className="default-align default-margin default-padding">
                <MuiTextField id={this.getID()} value={this.state.content} onChange={(event) => this.clickEvent(event)}
                              label={this.state.text} variant="outlined"
                              fullWidth/>
            </div>
        )
    }

    clickEvent(event) {
        let value = event.target.value
        let valueJson = JSON.stringify({text: value})
        window.fetch(this.state.updateRequest, {
            method: 'POST',
            mode: 'cors',
            cache: 'no-cache',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json',
            },
            redirect: 'follow',
            referrerPolicy: "no-referrer",
            body: valueJson
        });
        this.setState({content: value})
    }
}

export default TextField
