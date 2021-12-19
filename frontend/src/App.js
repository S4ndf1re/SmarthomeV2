import './css/App.css';
import React from 'react'
import ComponentViewer from './Components/ComponentViewer'
import ComponentPage from './Components/ComponentPage'
import {BrowserRouter as Router, Route,} from "react-router-dom"


class App extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            data: {
                containers: []
            }
        }
    }

    componentDidMount() {
        window.fetch("http://" + window.location.hostname + ":1337/gui", {
                credentials: "include"
            }
        ).then(response => response.json()).then(d => {
            this.setState({data: d})
        })
    }

    render() {
        return (
            <Router>
                {
                    this.state.data.containers.map(v =>
                        <Route key={v.name} exact path={"/" + v.name}>
                            <ComponentPage name={v.name} data={v.list} onInitRequest={v.onInitRequest}/>
                        </Route>
                    )
                }
                <Route exact path="/">
                    <ComponentViewer/>
                </Route>
            </Router>
        );
    }
}

export default App;
