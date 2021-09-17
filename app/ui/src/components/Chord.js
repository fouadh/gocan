import {useEffect, useRef} from "react";
import * as d3 from "d3";

export function Chord({data}) {
  const container = useRef(null);

  useEffect(() => {
    const width = 954;
    const radius = width / 2;

    let subscribe = true;

    const tree = d3.cluster()
      .size([2 * Math.PI, radius - 100]);

    const hierarchy = d3.hierarchy(data)
    .sort((a, b) => d3.ascending(a.height, b.height) ||
      d3.ascending(a.data.name, b.data.name));
    const root = tree(bilink(hierarchy));

    const svg = d3.select(container.current);
    svg.attr("viewBox", [-width / 2, -width / 2, width, width]);

    const degrees = root.leaves()
      .map(l => l.data)
      .filter(i => i.relations)
      .map(i => i.relations)
      .flatMap(i => i)
      .map(i => i.degree);

    const thicknessScale = d3.scaleLinear()
      .domain(degrees)
      .range([0, 10])

    const node = svg.append("g")
      .attr("font-family", "sans-serif")
      .attr("font-size", 10)
      .selectAll("g")
      .data(root.leaves())
      .join("g")
      .attr("transform", d => `rotate(${d.x * 180 / Math.PI - 90}) translate(${d.y},0)`)
      .append("text")
      .attr("dy", "0.31em")
      .attr("x", d => d.x < Math.PI ? 6 : -6)
      .attr("text-anchor", d => d.x < Math.PI ? "start" : "end")
      .attr("transform", d => d.x >= Math.PI ? "rotate(180)" : null)
      .text(d => d.data.name)
      .each(function (d) {
        d.text = this;
      })
      .on("mouseover", overed)
      .on("mouseout", outed)
      .call(text => text.append("title").text(d => `${id(d)}
${d.outgoing.length} outgoing
${d.incoming.length} incoming`));
    const colornone = "#ccc";
    const colorout = "#00f";
    const colorin = "#00f";
    const line = d3.lineRadial()
      .curve(d3.curveBundle.beta(0.85))
      .radius(d => d.y)
      .angle(d => d.x);

    const link = svg.append("g")
      .attr("stroke", colornone)
      .attr("fill", "none")
      .selectAll("path")
      .data(root.leaves().flatMap(leaf => leaf.outgoing))
      .join("path")
      .style("mix-blend-mode", "multiply")
      .attr("d", ([i, o]) => {
        console.log("incoming:", i, "outgoing:", o);
        return line(i.path(o));
      })
      .each(function (d) {
        d.path = this;
      });

    function bilink(root) {
      const map = new Map(root.leaves().map(d => [id(d), d]));
      for (const d of root.leaves()) {
        d.incoming = [];
        if (d.data.relations)
          d.outgoing = d.data.relations.map(i => [d, map.get(i.coupled)]);
        else
          d.outgoing = [];
      }
      for (const d of root.leaves()) for (const o of d.outgoing) {
        if (o[1] && o)
          o[1].incoming.push(o);
      }
      return root;
    }

    function id(node) {
      return `${node.parent ? id(node.parent) + "/" : ""}${node.data.name}`;
    }

    function overed(event, d) {
      link.style("mix-blend-mode", null);
      d3.select(this).attr("font-weight", "bold");
      d3.selectAll(d.incoming.map(d => d.path)).attr("stroke", colorin).raise();
      d3.selectAll(d.incoming.map(([d]) => d.text)).attr("fill", colorin).attr("font-weight", "bold");
      d3.selectAll(d.outgoing.map(d => d.path)).attr("stroke", colorout).raise();
      d3.selectAll(d.outgoing.map(([, d]) => d.text)).attr("fill", colorout).attr("font-weight", "bold");
    }

    function outed(event, d) {
      link.style("mix-blend-mode", "multiply");
      d3.select(this).attr("font-weight", null);
      d3.selectAll(d.incoming.map(d => d.path)).attr("stroke", null);
      d3.selectAll(d.incoming.map(([d]) => d.text)).attr("fill", null).attr("font-weight", null);
      d3.selectAll(d.outgoing.map(d => d.path)).attr("stroke", null);
      d3.selectAll(d.outgoing.map(([, d]) => d.text)).attr("fill", null).attr("font-weight", null);
    }

    return () => subscribe = false;
  }, [data]);

  return <svg ref={container}/>;
}