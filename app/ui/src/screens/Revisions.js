import {useEffect, useState} from "react";
import axios from "axios";
import {Timeline} from "../components/Timeline";
import {Spinner} from "../components/Spinner";
import {DateSelector} from "../components/DateSelector";

export function Revisions({sceneId, appId, date}) {
    const [dateRange, setDateRange] = useState(date);
    const [analyze, setAnalyze] = useState(true);
    const [revisions, setRevisions] = useState([]);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        let subscribed = true;
        if (analyze) {
            setLoading(true);
            let endpoint = `/api/scenes/${sceneId}/apps/${appId}/revisions`;
            if (dateRange.min) {
                if (dateRange.max) {
                    endpoint += `?after=${dateRange.min}&before=${dateRange.max}`
                } else {
                    endpoint += `?after=${dateRange.min}`
                }
            } else if (dateRange.max) {
                endpoint += `?before=${dateRange.max}`
            }
            axios.get(endpoint)
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
                }).finally(() => {
                    setAnalyze(false);
                    setLoading(false);
            });
        }

        return () => subscribed = false;
    }, [sceneId, appId, analyze, dateRange]);

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
        <DateSelector min={date.min} max={date.max} onChange={e => setDateRange(e)}/>
        <button onClick={e => setAnalyze(true)}>Submit</button>
        {screen}
    </>;
}