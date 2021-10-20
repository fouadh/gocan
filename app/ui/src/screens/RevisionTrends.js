import {useEffect, useState} from "react";
import axios from "axios";
import * as d3 from "d3";
import {MultiLineChart} from "../components/MultiLineChart";
import {Dropdown} from "primereact/dropdown";

export function RevisionTrends({sceneId, appId}) {
    const [trendName, setTrendName] = useState();
    const [trendNames, setTrendNames] = useState();
    const [trends, setTrends] = useState([]);
    const [transformations, setTransformations] = useState([]);

    useEffect(() => {
        let subscribed = true;
        axios.get(`/api/scenes/${sceneId}/apps/${appId}/revisions-trends`)
            .then(it => it.data)
            .then(it => it.trends)
            .then((it) => {
                if (subscribed) {
                    setTrendNames(it);
                }
            });
        return () => subscribed = false;
    }, [sceneId, appId])

    useEffect(() => {
        let subscribed = true;
        if (trendName) {
            axios.get(`/api/scenes/${sceneId}/apps/${appId}/revisions-trends/${trendName}`)
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
    }, [trendName, sceneId, appId]);

    let chart;
    if (trends && trends.length > 0) {
        chart = <MultiLineChart yLabel="Revisions Trends"
                                data={Object.values(trends)}
                                xAccessor={d => d3.timeParse('%Y-%m-%d')(d.x)}
                                yAccessor={d => d.y}
                                xFormatter={d3.timeFormat("%Y-%m-%d")}
                                legend={transformations}
        />;
    } else {
        chart = <></>;
    }

    return <div>
        <div>
            <label className="p-mr-2">Trend Name:</label>
            <Dropdown optionLabel="name"
                      optionValue="id"
                      options={trendNames}
                      placeholder="Select a trend"
                      value={trendName}
                      showClear={true}
                      onChange={(e) => {
                          setTrendName(e.value);
                      }}/>
        </div>
        { chart }
    </div>;
}