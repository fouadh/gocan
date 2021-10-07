import {ProgressSpinner} from "primereact/progressspinner";

export function Spinner() {
    return <>
        {<ProgressSpinner style={{top: "50%", position: "absolute", left: "50%", margin: "0", transform: "translate(-50%, -50%)"}}/>}
    </>;
}