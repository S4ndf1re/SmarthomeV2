const {Component} = require("react");


class Alert extends Component {
    constructor(props) {
        super(props);
        this.state = {
            message: props.message,
            name: props.name
        }
    }

    componentDidMount() {
        if (this.state.message !== null && this.state.message !== "") {
            alert(this.state.message)
        }
    }

    render() {
        return null
    }

}

export default Alert;