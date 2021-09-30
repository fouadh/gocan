import {useEffect, useState} from "react";
import axios from "axios";
import * as d3 from "d3";
import {MultiLineChart} from "../components/MultiLineChart";
import {Dropdown} from "primereact/dropdown";
import {Checkbox} from 'primereact/checkbox';
import './ComplexityAnalysis.css';

export function ComplexityAnalysis({sceneId, appId}) {
    const [complexity, setComplexity] = useState();
    const [displayedComplexities, setDisplayedComplexities] = useState([]);
    const [analysis, setAnalysis] = useState();
    const [analyses, setAnalyses] = useState([]);
    const [complexityTypes, setComplexityTypes] = useState(["indentations"]);

    const onComplexityTypeChange = (e) => {
        let selectedTypes = [...complexityTypes];
        if (e.checked) {
            selectedTypes.push(e.value);
        } else {
            if (selectedTypes.length > 1) {
                selectedTypes.splice(selectedTypes.indexOf(e.value), 1);
            }
        }
        setComplexityTypes(selectedTypes);
    }

    useEffect(() => {
        let subscribed = true;
        if (analysis) {
            axios.get(`/api/scenes/${sceneId}/apps/${appId}/complexity-analyses/${analysis}`)
                .then(it => it.data)
                .then(it => it.entries)
                .then(it => {
                    if (subscribed) {
                        const indentations = it.map((each) => ({x: each.date, y: each.indentations}));
                        const lines = it.map((each) => ({x: each.date, y: each.lines}));
                        const mean = it.map((each) => ({x: each.date, y: each.mean}));
                        const stdev = it.map((each) => ({x: each.date, y: each.stdev}));
                        const max = it.map((each) => ({x: each.date, y: each.max}));
                        const complexity = {indentations, lines, mean, stdev, max};
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

    useEffect(() => {
        if (!complexity) return;
        console.log({complexity});
        console.log({complexityTypes});

        const data = [];
        complexityTypes.forEach(type => {
            data.push(complexity[type]);
        });
        console.log({data});

        setDisplayedComplexities(data)
    }, [complexity, complexityTypes])

    let chart;
    let options;
    if (displayedComplexities && displayedComplexities.length > 0) {
        chart = <MultiLineChart label="Complexity"
                                data={displayedComplexities}
                                xAccessor={d => d3.isoParse(d.x)}
                                yAccessor={d => d.y}
                                xFormatter={d3.timeFormat("%Y-%m-%d")}
                                legend={complexityTypes}
        />;
        options = <>
            <div className="chart-options">
                <div className="p-col-12 chart-option">
                    <Checkbox inputId="cb1" value="lines" onChange={onComplexityTypeChange}
                              checked={complexityTypes.includes('lines')}></Checkbox>
                    <label htmlFor="cb1" className="p-checkbox-label">Lines</label>
                </div>
                <div className="p-col-12 chart-option">
                    <Checkbox inputId="cb2" value="indentations" onChange={onComplexityTypeChange}
                              checked={complexityTypes.includes('indentations')}></Checkbox>
                    <label htmlFor="cb2" className="p-checkbox-label">Indentations</label>
                </div>
                <div className="p-col-12 chart-option">
                    <Checkbox inputId="cb3" value="mean" onChange={onComplexityTypeChange}
                              checked={complexityTypes.includes('mean')}></Checkbox>
                    <label htmlFor="cb3" className="p-checkbox-label">Mean</label>
                </div>
                <div className="p-col-12 chart-option">
                    <Checkbox inputId="cb3" value="stdev" onChange={onComplexityTypeChange}
                              checked={complexityTypes.includes('stdev')}></Checkbox>
                    <label htmlFor="cb3" className="p-checkbox-label">Stdev</label>
                </div>
                <div className="p-col-12 chart-option">
                    <Checkbox inputId="cb3" value="max" onChange={onComplexityTypeChange}
                              checked={complexityTypes.includes('max')}></Checkbox>
                    <label htmlFor="cb3" className="p-checkbox-label">Max</label>
                </div>
            </div>
        </>
    } else {
        chart = <></>;
        options = <></>;
    }


    return (
        <>
            <div className="chart-form">
                <label className="p-mr-2">Select an analysis:</label>
                <Dropdown optionLabel="name"
                          optionValue="id"
                          options={analyses}
                          placeholder="Select an analysis"
                          value={analysis}
                          showClear={true}
                          onChange={(e) => {
                              setAnalysis(e.value);
                          }}/>
                {options}
            </div>
            {chart}
        </>
    )
        ;
}