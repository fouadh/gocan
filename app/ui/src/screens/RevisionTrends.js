import {useEffect, useState} from "react";
import axios from "axios";
import {Timeline} from "../components/Timeline";
import * as d3 from "d3";
import {BoundarySelector} from "./BoundarySelector";

export function RevisionTrends({sceneId, appId}) {
    const [boundary, setBoundary] = useState();
    const [trends, setTrends] = useState([]);
    const [transformations, setTransformations] = useState([]);

    useEffect(() => {
        let subscribed = true;
        if (boundary) {
            axios.get(`/api/scenes/${sceneId}/apps/${appId}/revisions-trends?boundary=${boundary}`)
                .then(it => it.data)
                .then(it => it.trends)
                .then(it => {
                    if (subscribed) {
                        console.log("it =>", it);
                        const transformations = [];
                        const data = [];

                        it.forEach(each => {
                            const row = {date: each.date};
                            each.revisions.forEach(rev => {
                                if (transformations.indexOf(rev.entity) < 0) {
                                    transformations.push(rev.entity);
                                }
                                row[rev.entity] = rev.numberOfRevisions;
                            });
                            data.push(row);
                        });
                        setTransformations(transformations);
                        setTrends(data);
                    }
                });
        } else {
            setTrends([]);
        }
        return () => subscribed = false;
    }, [boundary]);

    return <div>
        <BoundarySelector sceneId={sceneId} appId={appId} onChange={(e) => setBoundary(e.value)}/>
        {
            transformations.map((each) => {
                return <Timeline label={`Revisions for ${each}`} data={trends}
                                 xAccessor={(d) => d3.timeParse('%Y-%m-%d')(d.date)}
                                 yAccessor={(d) => d[each]} xFormatter={d3.timeFormat("%Y-%m-%d")}/>
            })}
    </div>;
}