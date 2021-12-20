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
        window.fetch("/gui", {
                credentials: "include",
                redirect: "follow"
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
                            <ComponentPage name={v.name} text={v.text} list={v.list} onInitRequest={v.onInitRequest}/>
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
