import Button from "./Button";
import Checkbox from "./Checkbox";
import TextField from "./TextField";
import Data from "./Data";
import Alert from "./Alert";

function getComponentForType(state) {
    if (state.type === "gui.Button") {
        return <Button key={state.name} name={state.name} text={state.text}
                       onClick={state.onClickRequest}/>
    } else if (state.type === "gui.Checkbox") {
        return <Checkbox key={state.name} name={state.name} text={state.text}
                         onStateClick={state.onOnStateRequest}
                         offStateClick={state.onOffStateRequest}
                         getStateClick={state.onGetStateRequest}/>
    } else if (state.type === "gui.TextField") {
        return <TextField key={state.name} name={state.name} text={state.text}
                          updateRequest={state.updateRequest}/>
    } else if (state.type === "gui.Data") {
        return <Data key={state.name} name={state.name} updateRequest={state.updateRequest}
                     updateSocket={state.updateSocket}/>
    } else if (state.type === "gui.Alert") {
        return <Alert key={state.name} name={state.name} message={state.message}/>
    }
}

export default getComponentForType