import React from 'react'
import Container from './Container'
import '../css/ComponentViewer.css'
import '../css/containerView.css'

class ComponentViewer extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            data: {
                containers: []
            }
        }
    }

    componentDidMount() {
        window.fetch("/gui", {
                credentials: "include",
                redirect: "follow"
            }
        ).then(response => {
            if (response.redirected) {
                window.location.href = "/login"
            } else {
                return response.json()
            }
        }).then(d => {
            this.setState({data: d})
        })
    }

    render() {
        const data = this.state.data
        return (
            <div className="max-div">
                <h1>Devices</h1>
                <div className="containerView blog-shadow-dreamy">
                    {
                        data["containers"].map(v => <Container name={v.name} text={v.text} list={v.list}/>)
                    }
                </div>
            </div>
        );
    }

}

export default ComponentViewer;