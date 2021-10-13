import {useEffect, useState} from "react";
import axios from "axios";
import {CirclePacking} from "../components/CirclePacking";
import {Spinner} from "../components/Spinner";

export function EntityCoupling({sceneId, appId}) {
    const [entity, setEntity] = useState("");
    const [error, setError] = useState();
    const [coupling, setCoupling] = useState();
    const [loading, setLoading] = useState(false);
    const [analyze, setAnalyze] = useState(false);

    useEffect(() => {
        let subscribed = true;
        if (analyze) {
            if (entity !== "") {
                setError(null);
                setLoading(true);
                axios.get(`/api/scenes/${sceneId}/apps/${appId}/entity-coupling?entity=${entity}`)
                    .then(it => it.data)
                    .then(it => {
                        if (subscribed) {
                            setCoupling(it);
                        }
                    })
                    .catch(() => setError("Unable to get coupling information"))
                    .finally(() => {
                    setLoading(false);
                    setAnalyze(false);
                });
            } else {
                setCoupling(null);
                setAnalyze(false);
            }
        }

        return () => subscribed = false;
    }, [sceneId, appId, analyze]);

    let screen;

    if (loading) {
        screen = <Spinner/>;
    } else if (coupling) {
        screen = <div style={{display: "flex", justifyContent: "center"}}>
            <CirclePacking width={975} height={975} data={coupling}/>
        </div>;
    } else if (error) {
        screen = <p>{error}</p>
    }

    return <>
        <div>
            <label htmlFor="entity">Entity</label>
            <input type="text" value={entity} onChange={e => setEntity(e.target.value)}/>
            <button onClick={e => setAnalyze(true)}>Submit</button>
        </div>
        {screen}
    </>
}