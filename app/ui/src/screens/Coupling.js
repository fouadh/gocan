import {useEffect, useState} from "react";
import axios from "axios";
import {Chord} from "../components/Chord";
import {Spinner} from "../components/Spinner";
import {DateSelector} from "../components/DateSelector";

export function Coupling({sceneId, appId, date}) {
    const [dateRange, setDateRange] = useState({date});
    const [analyze, setAnalyze] = useState(true);
    const [coupling, setCoupling] = useState();
    const [loading, setLoading] = useState(false);
    const [boundary] = useState("");

    useEffect(() => {
        let subscribe = true;
        if (analyze) {
            setLoading(true);
            let endpoint = `/api/scenes/${sceneId}/apps/${appId}/coupling-hierarchy`;
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
                .then(it => {
                    if (subscribe) {
                        setCoupling(it);
                    }
                })
                .finally(() => {
                    setLoading(false);
                    setAnalyze(false);
                });
        }

        return () => subscribe = false;
    }, [sceneId, appId, boundary, analyze, dateRange]);

    let screen;
    if (loading) {
        screen = <Spinner/>;
    } else if (coupling) {
        screen = <div style={{display: "flex", justifyContent: "center", height: "900px"}}>
            <Chord data={coupling}/>
        </div>
    } else {
        screen = <><p>No coupling found.</p></>
    }

    return <>
        <DateSelector min={date.min} max={date.max} onChange={e => setDateRange(e)}/>
        <button onClick={e => setAnalyze(true)}>Submit</button>
        {screen}
    </>;
}