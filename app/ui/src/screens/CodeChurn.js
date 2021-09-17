import {useEffect, useState} from "react";
import axios from "axios";
import * as d3 from "d3";
import {MultiLineChart} from "../components/MultiLineChart";

export function CodeChurn({sceneId, appId}) {
    const [codeChurn, setCodeChurn] = useState([]);

    useEffect(() => {
        let subscribed = true;
        axios.get(`/api/scenes/${sceneId}/apps/${appId}/code-churn`)
            .then(it => it.data)
            .then(it => it.codeChurn)
            .then(it => {
                if (subscribed) {
                    const added = it.map((each) => ({x: each.date, y: each.added}));
                    const deleted = it.map((each) => ({x: each.date, y: each.deleted}));
                    const churn = [added, deleted];
                    console.log({churn});
                    setCodeChurn(churn);
                }
            });

        return () => subscribed = false;
    }, [sceneId, appId]);

    return (
        <>
            <MultiLineChart label="Code Churn"
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