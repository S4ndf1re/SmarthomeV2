import '../css/Textfield.css'
import '../css/ComponentPage.css'
import '../css/shadow.css'
import '../css/index.css'
import React from 'react'

class TextField extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            updateRequest: props.updateRequest,
            name: props.name,
            text: props.text,
            checked: false
        }
    }

    getID() {
        return this.state.onUpdate + this.state.name
    }

    render() {
        return (
            <div className="textfieldClass default-margin default-padding blog-shadow-dreamy">
                <label className="defaultFont" htmlFor={this.getID()}>{this.state.text}</label>
                <input id={this.getID()} onChange={(key) => this.clickEvent(key)}
                       className="textInputClass defaultFont"/>
            </div>
        )
    }

    clickEvent(key) {
        console.log(key)
        let b = document.getElementById(this.getID()).value
        let c = JSON.stringify({text: b})
        window.fetch("http://" + window.location.hostname + ":1337/" + this.state.updateRequest, {
            method: 'POST',
            mode: 'cors',
            cache: 'no-cache',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json',
            },
            redirect: 'follow',
            referrerPolicy: "no-referrer",
            body: c
        });
    }
}

export default TextField
