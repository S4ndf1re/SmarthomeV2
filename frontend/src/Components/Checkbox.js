import '../css/Checkbox.css'
import '../css/ComponentPage.css'
import '../css/shadow.css'
import React from 'react'

class Checkbox extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            onStateClick: props.onStateClick,
            offStateClick: props.offStateClick,
            getStateClick: props.getStateClick,
            name: props.name,
            text: props.text,
            checked: false
        }
    }

    getID() {
        return this.state.onStateClick + this.state.offStateClick + this.state.name
    }

    render() {
        return (
            <div className="checkboxClass default-margin default-padding blog-shadow-dreamy">
                <input id={this.getID()} type="checkbox"
                       onChange={() => this.clickEvent()}
                       checked={this.state.checked}/>
                <label htmlFor={this.getID()}>{this.state.text}</label>
            </div>
        )
    }

    componentDidMount() {
        window.fetch("http://" + window.location.hostname + ":1337/" + this.state.getStateClick, {
            credentials: "include"
        }).then(data => data.json()).then(data => {
            document.getElementById(this.getID()).checked = data.status
            this.setState({checked: data.status})
        });
    }

    clickEvent() {
        let path = ""
        if (!this.state.checked) {
            path = this.state.onStateClick
        } else {
            path = this.state.offStateClick
        }
        window.fetch("http://" + window.location.hostname + ":1337/" + path, {
            credentials: "include"
        }).then(data => data.json()).then(data => {
            document.getElementById(this.getID()).checked = data.status
            this.setState({checked: data.status});
        });
    }
}

export default Checkbox
