import {useState} from "react";
import {Calendar} from 'primereact/calendar';

function formatDate(d) {
    const year = d.getFullYear();
    const month = d.getMonth() + 1;
    const day = d.getDate();

    const monthStr = month < 10 ? `0${month}` : `${month}`;
    const dayStr = day < 10 ? `0${day}` : `${day}`;

    return `${year}-${monthStr}-${dayStr}`;
}

export function DateSelector({min, max, onChange}) {
    const [dateRange] = useState({
        min: new Date(min + "T00:00"),
        max: new Date(max + "T00:00")
    });

    const notifyChange = (range) => {
        const date = {
            min: `${formatDate(range.min)}`,
            max: `${formatDate(range.max)}`,
        };
        onChange(date);
    };

    return <>
        <div className="p-field p-col-12 p-md-4 mr-4">
            <span className="p-float-label">
                <Calendar id="min" value={dateRange.min} onChange={e => {
                    notifyChange({...dateRange, min: e.value});
                }} dateFormat="yy-mm-dd"/>
                <label htmlFor="min">Min Date</label>
            </span>
        </div>
        <div className="p-field p-col-12 p-md-4 mr-4">
            <span className="p-float-label">
                <Calendar id="max" value={dateRange.max} onChange={e => {
                    notifyChange({...dateRange, max: e.value});
                }} dateFormat="yy-mm-dd"/>
                <label htmlFor="max">Max Date</label>
            </span>
        </div>
    </>
}