import '../css/Container.css'
import '../css/shadow.css'
import React from 'react'
import {Link} from "react-router-dom";


class Container extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            name: props.data.name,
            list: props.data.list
        }
    }


    render() {
        return (
            <Link to={this.state.name} className="container blog-shadow-dreamy">
                <p>{this.state.name}</p>
            </Link>
        );
    }
}

export default Container;
