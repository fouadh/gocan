import { useEffect, useRef } from 'react'
import * as d3 from 'd3'
import * as d3Cloud from 'd3-cloud'

export function WordCloud({ data, width = 600, height = 900 }) {
  const container = useRef(null)

  useEffect(() => {
    const svg = d3.select(container.current)
    svg.classed('word-cloud js-viz', true)
    const fontScale = d3.scaleLinear().range([20, 120])
    const fillScale = d3.scaleOrdinal(d3.schemeCategory10)
    const padding = 0
    const rotate = () => (~~(Math.random() * 6) - 3) * 30
    const minSize = d3.min(data, (d) => d.count) || 0
    const maxSize = d3.max(data, (d) => d.count) || 0
    fontScale.domain([minSize, maxSize])

    svg
      .attr('viewBox', `0 0 ${width} ${height}`)
      .attr('font-family', 'sans-serif')
      .attr('text-anchor', 'middle')

    const cloud = d3Cloud()
      .size([width, height])
      .words(data.map((d) => Object.create({ text: d.word, size: d.count })))
      .padding(padding)
      .rotate(rotate)
      .font('sans-serif')
      .fontSize((d) => fontScale(d.size || 0))
      .on('word', ({ size, x, y, rotate, text }) => {
        svg
          .append('text')
          .style('fill', fillScale(text || ''))
          .attr('font-size', size + 'px' || 0)
          .attr(
            'transform',
            (d) => 'translate(' + [x, y] + ')rotate(' + rotate + ')'
          )
          .text(text || '')
      })

    cloud.start()
  }, [data, width, height])

  return <svg ref={container} />
}
