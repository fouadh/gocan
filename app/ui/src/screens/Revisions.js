import {useEffect, useState} from "react";
import axios from "axios";
import {Timeline} from "../components/Timeline";
import {Spinner} from "../components/Spinner";

export function Revisions({sceneId, appId}) {
    const [revisions, setRevisions] = useState([]);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        let subscribed = true;
        setLoading(true);
        axios.get(`/api/scenes/${sceneId}/apps/${appId}/revisions`)
            .then(it => it.data)
            .then(it => it.revisions)
            .then(it => {
                const results = [];
                it.forEach((each, i) => results.push({x: i, y: each.numberOfRevisions}));
                return results;
            })
            .then(it => {
                if (subscribed) {
                    setRevisions(it);

                }
            }).finally(() => setLoading(false));

        return () => subscribed = false;
    }, [sceneId, appId]);

    let screen;
    if (loading) {
        screen = <Spinner/>;
    } else {
        screen = <Timeline data={revisions}
                           xAccessor={(d) => {
                               return d.x;
                           }}
                           yAccessor={(d) => d.y}
                           label="Revisions"
                           xFormatter={(d) => d.y}
        />;
    }

    return <>
        {screen}
    </>;
}