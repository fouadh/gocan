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
        console.log({date});
        onChange(date);
    };

    return <div>
        <label htmlFor="min">Min Date:</label>
        <Calendar value={dateRange.min} onChange={e => {
            notifyChange({...dateRange, min: e.value});
        }} dateFormat="yy-mm-dd"/>
        <label htmlFor="max">Max Date:</label>
        <Calendar value={dateRange.max} onChange={e => {
            notifyChange({...dateRange, max: e.value});
        }} dateFormat="yy-mm-dd"/>
    </div>
}