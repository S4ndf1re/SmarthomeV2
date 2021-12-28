import '../css/ComponentPage.css'
import '../css/shadow.css'
import React from 'react'
import MuiCheckbox from '@mui/material/Switch'
import {FormControlLabel, FormGroup} from "@mui/material";

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
            <div className="default-align default-margin default-padding">
                <FormGroup>
                    <FormControlLabel
                        control={<MuiCheckbox checked={this.state.checked} onChange={(evt) => this.clickEvent(evt)}/>}
                        label={this.state.text}/>
                </FormGroup>
            </div>
        )
    }

    componentDidMount() {
        window.fetch(this.state.getStateClick, {
            credentials: "include",
            redirect: "follow"
        }).then(data => data.json()).then(data => {
                this.setState({checked: data.status})
            }
        );
    }

    clickEvent(event) {
        let path = ""
        if (event.target.checked) {
            path = this.state.onStateClick
        } else {
            path = this.state.offStateClick
        }
        window.fetch(path, {
            credentials: "include",
            redirect: "follow"
        }).then(data => data.json()).then(data => {
            this.setState({checked: data.status});
        });
    }
}

export default Checkbox
