import '../css/Button.css'
import '../css/ComponentPage.css'
import '../css/shadow.css'
import React from 'react'

class Button extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            onClick: props.onClick,
            name: props.name,
            text: props.text
        }
    }

    render() {
        return (
            <button className="buttonClass default-margin blog-shadow-dreamy default-padding"
                    onClick={() => this.clickEvent()}>{this.state.text}</button>
        )
    }

    clickEvent() {
        window.fetch(this.state.onClick, {
            credentials: "include",
            redirect: "follow"
        })
    }
}

export default Button