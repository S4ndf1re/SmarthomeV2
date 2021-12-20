import '../css/Container.css'
import '../css/shadow.css'
import React from 'react'
import {Link} from "react-router-dom";


class Container extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            name: props.name,
            text: props.text,
            list: props.list
        }
    }


    render() {
        return (
            <Link to={this.state.name} className="container blog-shadow-dreamy">
                <p>{this.state.text}</p>
            </Link>
        );
    }
}

export default Container;
