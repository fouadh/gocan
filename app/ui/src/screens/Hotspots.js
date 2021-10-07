import {useEffect, useState} from "react";
import axios from "axios";
import {CirclePacking} from "../components/CirclePacking";
import {Spinner} from "../components/Spinner";

export function Hotspots({sceneId, appId}) {
    const [hospots, setHotspots] = useState();
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        let subscribed = true;
        setLoading(true);
        axios.get(`/api/scenes/${sceneId}/apps/${appId}/hotspots`)
            .then(it => it.data)
            .then(it => {
                if (subscribed) {
                    setHotspots(it);
                }
            }).finally(() => setLoading(false));
        ;

        return () => subscribed = false;
    }, [sceneId, appId]);

    let screen;

    if (loading) {
        screen = <Spinner/>;
    } else if (hospots) {
        screen = <div>
            <CirclePacking width={975} height={975} data={hospots}/>
        </div>;
    } else {
        screen = <><p>Unable to get hotspots</p></>
    }

    return <>{screen}</>
}