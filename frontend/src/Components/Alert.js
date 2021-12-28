import {Snackbar} from "@mui/material";
import MuiAlert from "@mui/material/Alert"
import * as React from 'react'


class Alert
    extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            message: props.message,
            name: props.name,
            severity: props.severity,
            open: true
        }
    }

    componentDidMount() {
    }

    handleClose(event, reason) {
        if (reason === 'clickaway') {
            return
        }
        this.setState({open: false})
    }

    render() {
        const {vertical, horizontal} = {vertical: 'top', horizontal: 'right'}
        if (this.state.message !== null && this.state.message !== "") {
            return (
                <Snackbar open={this.state.open} autoHideDuration={6000}
                          onClose={(event, reason) => this.handleClose(event, reason)}
                          anchorOrigin={{vertical, horizontal}}>
                    <MuiAlert severity={this.state.severity}
                              onClose={(event, reason) => this.handleClose(event, reason)}
                              sx={{width: '100%'}}>
                        {this.state.message}
                    </MuiAlert>
                </Snackbar>
            )
        }
        return null
    }

}

export default Alert;