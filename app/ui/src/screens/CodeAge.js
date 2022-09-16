import {useEffect, useState} from "react";
import axios from "axios";
import {CirclePacking} from "../components/CirclePacking";
import {Spinner} from "../components/Spinner";
import {DateSelector} from "../components/DateSelector";
import {Button} from 'primereact/button';
import {Calendar} from "primereact/calendar";

function formatDate(d) {
    const year = d.getFullYear();
    const month = d.getMonth() + 1;
    const day = d.getDate();

    const monthStr = month < 10 ? `0${month}` : `${month}`;
    const dayStr = day < 10 ? `0${day}` : `${day}`;

    return `${year}-${monthStr}-${dayStr}`;
}


export function CodeAge({sceneId, appId, date}) {
    const [dateRange, setDateRange] = useState(date);
    const [initialDate, setInitialDate] = useState(new Date(date.max + "T00:00"));
    const [analyze, setAnalyze] = useState(true);
    const [hospots, setHotspots] = useState();
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        let subscribed = true;
        if (analyze) {
            setLoading(true);
            let endpoint = `/api/scenes/${sceneId}/apps/${appId}/code-age`;
            let params = new URLSearchParams();
            if (dateRange.min) {
                params.append("after", dateRange.min);
            }
            if (dateRange.max) {
                params.append("before", dateRange.max);
            }
            params.append("initialDate", formatDate(initialDate));
            axios.get(`${endpoint}?${params}`)
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
    }, [sceneId, appId, analyze, dateRange, initialDate]);

    let screen;

    if (loading) {
        screen = <Spinner/>;
    } else if (hospots) {
        screen = <div style={{display: "flex", justifyContent: "center"}}>
            <CirclePacking width={975} height={975} data={hospots}/>
        </div>;
    } else {
        screen = <><p>Unable to get code age hotspots</p></>
    }

    return <>
        <div className="card mt-4">
            <div className="flex align-items-center">
                <DateSelector min={date.min} max={date.max} onChange={e => setDateRange(e)}/>
                <div className="p-field p-col-12 p-md-4 mr-4">
                <span className="p-float-label">
                    <Calendar id="initialDate" value={initialDate} onChange={e => {
                        setInitialDate(e.value);
                    }} dateFormat="yy-mm-dd"/>
                    <label htmlFor="initialDate">Counting From</label>
                </span>
                    </div>

                <Button label="Submit" onClick={e => setAnalyze(true)}/>
            </div>
        </div>
        {screen}
    </>
}