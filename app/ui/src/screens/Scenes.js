import {useEffect, useState} from "react";
import axios from "axios";
import {Link} from 'react-router-dom';

function getScenes() {
    return axios
        .get('/api/scenes')
        .then(it => {
            return it.data;
        })
        .then(it => it.scenes);
}

export function Scenes() {
    const [scenes, setScenes] = useState([]);
    useEffect(() => {
        let subscribed = true;
        getScenes()
            .then(it => {
                if (subscribed && it)
                    setScenes(it);
            });

        return (() => subscribed = false);
    }, []);

    let screen;
    if (scenes && scenes.length > 0) {
        screen = <div className="scenes ml-4 mt-4">
            {
                scenes.map((each) => <div key={each.name}>
                    <Link to={`/scenes/${each.id}`}>
                        <span data-name={each.name}>
                            {each.name}
                        </span>
                    </Link>
                </div>)
            }
        </div>
    } else {
        screen = <>
            No scene found. Please use the command <strong>gocan create-scene</strong> to create some.
        </>
    }

    return <>{screen}</>;
}