import * as d3 from 'd3'
import React, { useEffect, useRef, useState } from 'react'

const color = d3
  .scaleLinear()
  .domain([-1, 5])
  .range(['hsl(185,60%,99%)', 'hsl(187,40%,70%)'])
  .interpolate(d3.interpolateHcl)

function defaultFillColor(d) {
  return (d.data.weight || 0) > 0.0
    ? 'darkred'
    : d.children
    ? color(d.depth)
    : 'WhiteSmoke'
}

function defaultFillOpacity(d) {
  return d.data.weight || 1
}

function defaultTitle(d) {
  const ancestors = d.ancestors()
  const hierarchy = [...ancestors.slice(0, ancestors.length - 1)]
  return `${hierarchy
    .map((d) => d.data.name)
    .reverse()
    .join('/')}\n${d3.format(',d')(d.value || 0)}`
}

export function CirclePacking({
  data = { name: 'empty' },
  width = 400,
  height = 400,
  zoomToPath,
  fillColor = defaultFillColor,
  fillOpacity = defaultFillOpacity,
  setTitle = defaultTitle,
  onNodeClicked = () => {},
}) {
  const container = useRef(null)
  const [svg, setSvg] = useState()
  const [label, setLabel] = useState()
  const [node, setNode] = useState()
  const [focus, setFocus] = useState()
  const [view, setView] = useState()

  function zoom(d) {
    let viewElement = view
    if (!viewElement) {
      viewElement = [d.x, d.y, d.r * 2]
      setView(viewElement)
    }
    const interpolator = d3.interpolateZoom(viewElement, [
      focus.x,
      focus.y,
      focus.r * 2,
    ])
    const transition = svg
      .transition()
      .duration(interpolator.duration)
      .tween('zoom', () => {
        return (t) => {
          zoomTo(interpolator(t))
        }
      })

    label
      .filter(function (d) {
        return d.parent === focus || this.style.display === 'inline'
      })
      .transition(transition)
      .style('fill-opacity', (d) => (d.parent === focus ? 1 : 0))
      .on('start', function (d) {
        if (d.parent === focus) this.style.display = 'inline'
      })
      .on('end', function (d) {
        if (d.parent !== focus) this.style.display = 'none'
      })
  }

  function zoomTo(v) {
    setView(v)
    const k = width / v[2]
    label.attr(
      'transform',
      (d) => `translate(${(d.x - v[0]) * k},${(d.y - v[1]) * k})`
    )
    node.attr(
      'transform',
      (d) => `translate(${(d.x - v[0]) * k},${(d.y - v[1]) * k})`
    )
    node.attr('r', (d) => d.r * k)
  }

  useEffect(() => {
    if (focus) {
      zoom(focus)
    }
  }, [focus, svg, node, label])

  useEffect(() => {
    const color = d3
      .scaleLinear()
      .domain([-1, 5])
      .range(['hsl(185,60%,99%)', 'hsl(187,40%,70%)'])
      .interpolate(d3.interpolateHcl)
    const pack = d3.pack().size([width, height]).padding(3)
    const root = pack(
      d3
        .hierarchy(data)
        .sum((d) => d.size || 0)
        .sort((a, b) => (b.value || 0) - (a.value || 0))
    )
    let svgElement = buildSvg()
    setSvg(svgElement)
    const nodeElement = buildNodes(svgElement)
    setNode(nodeElement)
    buildTitle(nodeElement)
    const labelElement = buildLabels(svgElement)
    setLabel(labelElement)
    if (zoomToPath) {
      let nodes = root.children
      let elements = zoomToPath.split('/')
      let node
      while (elements.length > 0) {
        const element = elements.pop()
        let found = false
        for (let i = 0; i < nodes.length && !found; i++) {
          if (nodes[i].data.name === element) {
            found = true
            node = nodes[i]
            nodes = node.children
          }
        }
      }
      setFocus(node)
    } else {
      setFocus(root)
    }

    function buildSvg() {
      const svg = d3.select(container.current)
      svg.classed('circle-packing js-viz', true)
      svg
        .style('display', 'block')
        .style('cursor', 'pointer')
        .attr('width', width)
        .attr('height', height)
        .attr('viewBox', `-${width / 2} -${height / 2} ${width} ${height}`)
        .style('margin', '0 -14px')
        .attr('background', color(0))
      return svg
    }

    function buildNodes(svg) {
      return svg
        .append('g')
        .selectAll('circle')
        .data(root.descendants())
        .join('circle')
        .style('fill', fillColor)
        .style('fill-opacity', fillOpacity)
        .on('mouseover', function () {
          d3.select(this).attr('stroke', '#000000')
        })
        .on('mouseout', function () {
          d3.select(this).attr('stroke', null)
        })
        .on('click', (event, d) => {
          onNodeClicked(d)
          return (
            d !== focus && d.children && (setFocus(d), event.stopPropagation())
          )
        })
        .attr('role', (d) => (focus !== d ? 'link' : null))
    }

    function buildLabels(svg) {
      return svg
        .style('font', '10px sans-serif')
        .attr('text-anchor', 'middle')
        .append('g')
        .selectAll('text')
        .data(root.descendants())
        .join('text')
        .style('fill-opacity', (d) => (d.parent === root ? 1 : 0))
        .style('display', (d) => (d.parent === root ? 'inline' : 'none'))
        .text((d) => d.data.name)
    }

    function buildTitle(node) {
      node.append('title').text(setTitle)
    }
  }, [height, width, data, fillColor, fillOpacity, setTitle, zoomToPath])

  return <svg ref={container} />
}
