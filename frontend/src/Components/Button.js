import React from 'react'
import MuiButton from '@mui/material/Button';
import '../css/shadow.css'
import '../css/ComponentPage.css'

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
            <div className="default-align default-margin default-padding">
                <MuiButton variant="outlined" onClick={() => this.clickEvent()}>{this.state.text}</MuiButton>
            </div>
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