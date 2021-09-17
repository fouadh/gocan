import {useEffect, useState} from "react";
import axios from "axios";
import {WordCloud} from "../components/WordCloud";

export function ModusOperandi({sceneId, appId}) {
  const [messages, setMessages] = useState([]);

  useEffect(() => {
    let subscribed = true;
    axios.get(`/api/scenes/${sceneId}/apps/${appId}/modus-operandi`)
      .then(it => it.data)
      .then(it => it.modusOperandi)
      .then(it => {
        if (subscribed)
          setMessages(it);
      });

    return () => subscribed = false;
  }, [sceneId, appId]);


  return <WordCloud data={messages.slice(0, 100)}/>;
}