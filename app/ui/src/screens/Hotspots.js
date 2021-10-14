import {useEffect, useState} from "react";
import axios from "axios";
import {CirclePacking} from "../components/CirclePacking";
import {Spinner} from "../components/Spinner";
import {DateSelector} from "../components/DateSelector";

export function Hotspots({sceneId, appId}) {
    const [dateRange, setDateRange] = useState({});
    const [analyze, setAnalyze] = useState(true);
    const [hospots, setHotspots] = useState();
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        let subscribed = true;
        if (analyze) {
            setLoading(true);
            let endpoint = `/api/scenes/${sceneId}/apps/${appId}/hotspots`;
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
                    if (subscribed) {
                        setHotspots(it);
                    }
                }).finally(() => {
                setLoading(false);
                setAnalyze(false);
            });
        }

        return () => subscribed = false;
    }, [sceneId, appId, analyze]);

    let screen;

    if (loading) {
        screen = <Spinner/>;
    } else if (hospots) {
        screen = <div style={{display: "flex", justifyContent: "center"}}>
            <CirclePacking width={975} height={975} data={hospots}/>
        </div>;
    } else {
        screen = <><p>Unable to get hotspots</p></>
    }

    return <>
        <DateSelector onChange={e => setDateRange(e)}/>
        <button onClick={e => setAnalyze(true)}>Submit</button>
        {screen}
    </>
}