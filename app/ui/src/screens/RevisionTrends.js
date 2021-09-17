import {useEffect, useState} from "react";
import axios from "axios";
import * as d3 from "d3";
import {BoundarySelector} from "./BoundarySelector";
import {MultiLineChart} from "../components/MultiLineChart";

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
                        const map = {};

                        it.forEach(each => {
                           const date = each.date;
                           const revs = each.revisions;

                           revs.forEach(rev => {
                              if (!map[rev.entity]) {
                                  map[rev.entity] = [];
                              }

                              map[rev.entity].push({ x: date, y: rev.numberOfRevisions });
                           });
                        });

                        let keys = Object.keys(map);
                        setTransformations(keys);
                        const trends = keys.map(t => map[t]);
                        setTrends(trends);
                    }
                });
        } else {
            setTrends({});
        }
        return () => subscribed = false;
    }, [boundary, sceneId, appId]);

    return <div>
        <BoundarySelector sceneId={sceneId} appId={appId} onChange={(e) => setBoundary(e.value)}/>
        <MultiLineChart label="Revisions Trends"
                        data={Object.values(trends)}
                        xAccessor={d => d3.timeParse('%Y-%m-%d')(d.x)}
                        yAccessor={d => d.y}
                        xFormatter={d3.timeFormat("%Y-%m-%d")}
                        legend={transformations}
        />
    </div>;
}