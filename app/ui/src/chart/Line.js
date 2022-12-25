import React from 'react'
import * as d3 from 'd3'

const Line = ({
  type,
  data,
  xAccessor,
  yAccessor,
  y0Accessor,
  interpolation,
  ...props
}) => {
  const lineGenerator = d3[type]()
    .x(xAccessor)
    .y(yAccessor)
    .curve(interpolation)

  if (type === 'area') {
    lineGenerator.y0(y0Accessor).y1(yAccessor)
  }

  return (
    <path
      stroke="#9980FA"
      className={`Line Line--type-${type}`}
      d={lineGenerator(data)}
      {...props}
    />
  )
}

Line.defaultProps = {
  type: 'line',
  y0Accessor: 0,
  interpolation: d3.curveMonotoneX,
}

export { Line }
