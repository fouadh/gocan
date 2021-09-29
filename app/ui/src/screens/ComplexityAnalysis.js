import {useEffect, useState} from "react";
import axios from "axios";
import * as d3 from "d3";
import {MultiLineChart} from "../components/MultiLineChart";
import {Dropdown} from "primereact/dropdown";

export function ComplexityAnalysis({sceneId, appId}) {
    const [complexity, setComplexity] = useState([]);
    const [analysis, setAnalysis] = useState();
    const [analyses, setAnalyses] = useState([]);

    useEffect(() => {
        let subscribed = true;
        if (analysis) {
            axios.get(`/api/scenes/${sceneId}/apps/${appId}/complexity-analyses/${analysis}`)
                .then(it => it.data)
                .then(it => it.entries)
                .then(it => {
                    if (subscribed) {
                        const indentations = it.map((each) => ({x: each.date, y: each.indentations}));
                        const complexity = [indentations];
                        setComplexity(complexity);
                    }
                });
        }

        return () => subscribed = false;
    }, [sceneId, appId, analysis]);

    useEffect(() => {
        let subscribed = true;
        axios.get(`/api/scenes/${sceneId}/apps/${appId}/complexity-analyses`)
            .then(it => it.data)
            .then(it => it.analyses)
            .then(it => {
               if (subscribed) {
                   setAnalyses(it);
               }
            });
        return () => subscribed = false;
    }, []);

    return (
        <>
            <label className="p-mr-2">Analysis:</label>
            <Dropdown optionLabel="name"
                      optionValue="id"
                      options={analyses}
                      placeholder="Select an analysis"
                      value={analysis}
                      showClear={true}
                      onChange={(e) => {
                          setAnalysis(e.value);
                      }}/>

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