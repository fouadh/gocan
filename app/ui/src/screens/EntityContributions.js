import {useChartDimensions} from "../hooks/useChartDimensions";
import Chart from "../chart/Chart";
import * as d3 from "d3"
import {DateSelector} from "../components/DateSelector";
import {Button} from "primereact/button";
import {useEffect, useState} from "react";
import {Spinner} from "../components/Spinner";
import axios from "axios";
import {Autocomplete} from "../components/Autocomplete";

export function EntityContributions({sceneId, appId, date, entities}) {
    const [ref, dimensions] = useChartDimensions({height: "800px", width: "500px"});
    const [entity, setEntity] = useState("");
    const [dateRange, setDateRange] = useState(date);
    const [analyze, setAnalyze] = useState(true);
    const [error, setError] = useState();
    const [loading, setLoading] = useState(false);
    const [rectangles, setRectangles] = useState([]);
    const [authors, setAuthors] = useState([]);

    useEffect(() => {
        let subscribed = true;
        const colorScale = d3.scaleOrdinal(d3.schemeSet1);

        if (analyze) {
            if (entity !== "") {
                setError(null);
                setLoading(true);
                let endpoint = `/api/scenes/${sceneId}/apps/${appId}/entity-contributions?entity=${entity}`;
                if (dateRange.min) {
                    if (dateRange.max) {
                        endpoint += `&after=${dateRange.min}&before=${dateRange.max}`
                    } else {
                        endpoint += `&after=${dateRange.min}`
                    }
                } else if (dateRange.max) {
                    endpoint += `&before=${dateRange.max}`
                }

                axios.get(endpoint)
                    .then(it => it.data)
                    .then(it => it.contributions)
                    .then(it => {
                        if (subscribed) {
                            setRectangles(buildRectangles(it));
                            setAuthors(it.map((c, index) => ({
                                name: c.dev,
                                fill: colorScale(index + 1)
                            })));
                        }
                    })
                    .catch(() => setError("Unable to get contributions information"))
                    .finally(() => {
                        setLoading(false);
                        setAnalyze(false);
                    });
            } else {
                setAnalyze(false);
                setRectangles([]);
                setAuthors([]);
            }
        }

        const buildRectangles = (contributions) => {
            const totalContributions = contributions.map(c => c.contributions).reduce((a, c) => a + c, 0);

            const initialRectangle = {id: 0, x: 10, y: 10, height: 600, width: 400, fill: "#69b3a2"};
            const totalArea = initialRectangle.height * initialRectangle.width;
            const rectangle0 = {
                id: 1,
                x: initialRectangle.x,
                y: initialRectangle.y,
                height: initialRectangle.height,
                width: (contributions[0].contributions / totalContributions) * totalArea / initialRectangle.height,
                fill: colorScale(1)
            };

            const rectangles = [initialRectangle, rectangle0];

            for (let i = 2; i <= contributions.length; i++) {
                if (i % 2 === 0) {
                    rectangles.push({
                        id: i,
                        x: rectangles[i - 1].x + rectangles[i - 1].width,
                        y: rectangles[i - 1].y,
                        height: (contributions[i - 1].contributions / totalContributions) * totalArea / (rectangles[i - 2].width - rectangles[i - 1].width),
                        width: rectangles[i - 2].width - rectangles[i - 1].width,
                        fill: colorScale(i)
                    });
                } else {
                    rectangles.push({
                        id: i,
                        x: rectangles[i - 1].x,
                        y: rectangles[i - 1].y + rectangles[i - 1].height,
                        height: rectangles[i - 2].height - rectangles[i - 1].height,
                        width: (contributions[i - 1].contributions / totalContributions) * totalArea / (rectangles[i - 2].height - rectangles[i - 1].height),
                        fill: colorScale(i)
                    });
                }
            }

            return rectangles;
        }

        return () => subscribed = false;
    }, [sceneId, appId, analyze, dateRange, entity]);


    return <>
        {loading && <Spinner/>}
        {error && <p>{error}</p>}
        <div className="card mt-4">
            <div className="flex align-items-start">
                <div className="p-field p-col-12 p-md-4 mr-4">
                    <span className="p-float-label autocomplete">
                        <Autocomplete suggestions={entities} onChange={e => setEntity(e.value)}/>
                        <label htmlFor="entity">Entity</label>
                    </span>
                </div>
                <DateSelector min={date.min} max={date.max} onChange={e => setDateRange(e)}/>
                <Button label="Submit" onClick={e => setAnalyze(true)}/>
            </div>
            <div style={{display: "flex"}}>
                <div className="p-mr-5">
                    <ul style={{
                        width: "250px",
                        marginTop: "100px"
                    }}>
                        {authors.map(a =>
                            (
                                <li key={a.name} style={{listStyleType: 'none'}}>
                <span
                    style={{
                        display: 'inline-block',
                        width: '20px',
                        height: '20px',
                        marginRight: '.5em',
                        backgroundColor: a.fill
                    }}
                ></span>
                                    <span>{a.name}</span>
                                </li>
                            ))}
                    </ul>
                </div>
                <div className="js-viz" ref={ref}>
                    <Chart dimensions={dimensions}>
                        {
                            rectangles.map(r => <rect stroke={"#eee"} x={r.x} y={r.y} height={r.height} width={r.width}
                                                      fill={r.fill} key={r.id}/>)
                        }
                    </Chart>
                </div>
            </div>
        </div>
    </>
}