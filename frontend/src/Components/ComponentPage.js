import React from 'react'
import '../css/shadow.css'
import '../css/ComponentPage.css'
import '../css/containerView.css'
import getComponentForType from "./Util";

class ComponentPage extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            name: props.name,
            list: props.data,
            onInitRequest: props.onInitRequest
        }
        this.components = []
        this.updateList()
    }

    updateList() {
        this.state.list.forEach(state => {
            this.components.push(getComponentForType(state))
        })
    }

    componentDidMount() {
        this.updateList()
        window.fetch("http://" + window.location.hostname + ":1337/" + this.state.onInitRequest, {
                credentials: "include"
            }
        ).catch(err => console.log(err))
    }

    render() {
        return (
            <div>
                <h1> {this.state.name}
                </h1>
                <div className="containerView blog-shadow-dreamy">
                    <>
                        {this.components}
                    </>
                </div>
            </div>
        )
    }

}

export default ComponentPage