import {useEffect, useState} from "react";
import axios from "axios";
import * as d3 from "d3";
import {MultiLineChart} from "../components/MultiLineChart";

export function ComplexityAnalysis({sceneId, appId}) {
    const [complexity, setComplexity] = useState([]);

    useEffect(() => {
        let subscribed = true;
        axios.get(`/api/scenes/${sceneId}/apps/${appId}/complexity-analysis/analysis-1`)
            .then(it => it.data)
            .then(it => it.entries)
            .then(it => {
                if (subscribed) {
                    const indentations = it.map((each) => ({x: each.date, y: each.indentations}));
                    const complexity = [indentations];
                    setComplexity(complexity);
                }
            });

        return () => subscribed = false;
    }, [sceneId, appId]);

    return (
        <>
            <MultiLineChart label="Complexity"
                            data={complexity}
                            xAccessor={d => d3.isoParse(d.x)}
                            yAccessor={d => d.y}
                            xFormatter={d3.timeFormat("%Y-%m-%d")}
                            legend={["Complexity"]}
            />
        </>
    )
        ;
}