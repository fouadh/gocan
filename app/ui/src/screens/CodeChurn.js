import {useEffect, useState} from "react";
import axios from "axios";
import * as d3 from "d3";
import {MultiLineChart} from "../components/MultiLineChart";
import {DateSelector} from "../components/DateSelector";

export function CodeChurn({sceneId, appId, date}) {
    const [dateRange, setDateRange] = useState(date);
    const [analyze, setAnalyze] = useState(true);
    const [codeChurn, setCodeChurn] = useState([]);

    useEffect(() => {
        let subscribed = true;

        if (analyze) {
            let endpoint = `/api/scenes/${sceneId}/apps/${appId}/code-churn`;
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
                .then(it => it.codeChurn)
                .then(it => {
                    if (subscribed) {
                        const added = it.map((each) => ({x: each.date, y: each.added}));
                        const deleted = it.map((each) => ({x: each.date, y: each.deleted}));
                        const churn = [added, deleted];
                        setCodeChurn(churn);
                    }
                })
                .finally(() => {
                    setAnalyze(false);
                });
        }

        return () => subscribed = false;
    }, [sceneId, appId, analyze, dateRange]);

    return (
        <>
            <DateSelector min={date.min} max={date.max} onChange={e => setDateRange(e)}/>
            <button onClick={e => setAnalyze(true)}>Submit</button>
            <MultiLineChart yLabel="Code Churn"
                            data={codeChurn}
                            xAccessor={d => d3.timeParse('%Y-%m-%d')(d.x)}
                            yAccessor={d => d.y}
                            xFormatter={d3.timeFormat("%Y-%m-%d")}
                            legend={["Added", "Deleted"]}
            />
        </>
    )
        ;
}