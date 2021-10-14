import {useEffect, useState} from "react";

export function DateSelector({min, max, onChange}) {
    const [dateRange, setDateRange] = useState({
        min: min,
        max: max
    });

    useEffect(() => {
        onChange(dateRange);
    }, [dateRange]);

    return <div>
        <label htmlFor="min">Min Date:</label>
        <input type="text" value={dateRange.min} onChange={e => setDateRange({...dateRange, min: e.target.value})}/>
        <label htmlFor="max">Max Date:</label>
        <input type="text" value={dateRange.max} onChange={e => setDateRange({...dateRange, max: e.target.value})}/>
    </div>
}