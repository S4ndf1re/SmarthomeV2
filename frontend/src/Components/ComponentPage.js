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
            text: props.text,
            list: props.list,
            onInitRequest: props.onInitRequest,
            onUnloadRequest: props.onUnloadRequest
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
        window.fetch(this.state.onInitRequest, {
                credentials: "include",
                redirect: "follow"
            }
        ).catch(err => console.log(err))
    }

    componentWillUnmount() {
        window.fetch(this.state.onUnloadRequest, {
                credentials: "include",
                redirect: "follow"
            }
        ).catch(err => console.log(err))
    }

    render() {
        return (
            <div>
                <h1> {this.state.text}
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