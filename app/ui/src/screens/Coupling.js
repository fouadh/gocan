import {useEffect, useState} from "react";
import axios from "axios";
import {Chord} from "../components/Chord";

export function Coupling({sceneId, appId}) {
    const [coupling, setCoupling] = useState();
    const [boundary] = useState("");

    useEffect(() => {
        let subscribe = true;
        axios.get(`/api/scenes/${sceneId}/apps/${appId}/coupling-hierarchy`)
            .then(it => it.data)
            .then(it => {
                if (subscribe) {
                    setCoupling(it);
                }
            });

        return () => subscribe = false;
    }, [sceneId, appId, boundary]);

    return <div>
        {coupling ? <Chord data={coupling}/> : null}
    </div>;
}