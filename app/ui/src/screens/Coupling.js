import {useEffect, useState} from "react";
import axios from "axios";
import {Chord} from "../components/Chord";
import {Spinner} from "../components/Spinner";

export function Coupling({sceneId, appId}) {
    const [coupling, setCoupling] = useState();
    const [loading, setLoading] = useState(false);
    const [boundary] = useState("");

    useEffect(() => {
        let subscribe = true;
        setLoading(true);
        axios.get(`/api/scenes/${sceneId}/apps/${appId}/coupling-hierarchy`)
            .then(it => it.data)
            .then(it => {
                if (subscribe) {
                    setCoupling(it);
                }
            })
            .finally(() => setLoading(false));

        return () => subscribe = false;
    }, [sceneId, appId, boundary]);

    let screen;
    if (loading) {
        screen = <Spinner/>;
    } else if (coupling) {
        screen = <Chord data={coupling}/>
    } else {
        screen = <><p>No coupling found.</p></>
    }

    return <>{screen}</>;
}