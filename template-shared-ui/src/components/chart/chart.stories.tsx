import type { Meta, StoryObj } from '@storybook/react'
import { Area, AreaChart, Bar, BarChart } from 'recharts'
import { CartesianGrid, Line, LineChart, XAxis, YAxis } from 'recharts'
import { ChartContainer, ChartTooltip, ChartTooltipContent } from './chart'

const meta: Meta<typeof ChartContainer> = {
  title: 'Visualizations/Chart',
  component: ChartContainer,
}

export default meta
type Story = StoryObj<typeof meta>

const data = [
  {
    name: 'Jan',
    total: 167,
    revenue: 234,
    profit: 389,
    sales: 430,
    growth: 120,
  },
  {
    name: 'Feb',
    total: 245,
    revenue: 278,
    profit: 420,
    sales: 380,
    growth: 180,
  },
  {
    name: 'Mar',
    total: 321,
    revenue: 189,
    profit: 490,
    sales: 480,
    growth: 230,
  },
  {
    name: 'Apr',
    total: 356,
    revenue: 239,
    profit: 520,
    sales: 520,
    growth: 280,
  },
  {
    name: 'May',
    total: 270,
    revenue: 349,
    profit: 450,
    sales: 600,
    growth: 320,
  },
  {
    name: 'Jun',
    total: 429,
    revenue: 319,
    profit: 580,
    sales: 650,
    growth: 340,
  },
]

const chartConfig = {
  total: {
    theme: {
      light: 'var(--chart-1)',
      dark: 'var(--chart-1)',
    },
    label: 'Total',
  },
  revenue: {
    theme: {
      light: 'var(--chart-2)',
      dark: 'var(--chart-2)',
    },
    label: 'Revenue',
  },
  profit: {
    theme: {
      light: 'var(--chart-3)',
      dark: 'var(--chart-3)',
    },
    label: 'Profit',
  },
  sales: {
    theme: {
      light: 'var(--chart-4)',
      dark: 'var(--chart-4)',
    },
    label: 'Sales',
  },
  growth: {
    theme: {
      light: 'var(--chart-5)',
      dark: 'var(--chart-5)',
    },
    label: 'Growth',
  },
}

export const LineChartExample: Story = {
  render: () => (
    <div className="h-[400px] w-[800px]">
      <ChartContainer config={chartConfig}>
        <LineChart data={data}>
          <CartesianGrid strokeDasharray="3 3" className="stroke-border" />
          <XAxis dataKey="name" stroke="#888888" fontSize={12} tickLine={false} axisLine={false} />
          <YAxis
            stroke="#888888"
            fontSize={12}
            tickLine={false}
            axisLine={false}
            tickFormatter={(value) => `$${value}`}
          />
          <ChartTooltip
            content={({ active, payload }) => {
              if (!active || !payload) return null
              return (
                <ChartTooltipContent>
                  {payload.map((item) => (
                    <div key={item.dataKey} className="flex items-center justify-between gap-2">
                      <span className="text-muted-foreground">{item.name}</span>
                      <span className="font-bold">${item.value}</span>
                    </div>
                  ))}
                </ChartTooltipContent>
              )
            }}
          />
          <Line
            type="monotone"
            dataKey="total"
            strokeWidth={2}
            activeDot={% raw %}{{
              r: 6,
              style: { fill: 'var(--color-total)', opacity: 0.8 },
            }}{% endraw %}
            style={% raw %}{{
              stroke: 'var(--color-total)',
              opacity: 0.8,
            }}{% endraw %}
          />
          <Line
            type="monotone"
            dataKey="revenue"
            strokeWidth={2}
            activeDot={% raw %}{{
              r: 6,
              style: { fill: 'var(--color-revenue)', opacity: 0.8 },
            }}{% endraw %}
            style={% raw %}{{
              stroke: 'var(--color-revenue)',
              opacity: 0.8,
            }}{% endraw %}
          />
        </LineChart>
      </ChartContainer>
    </div>
  ),
}

export const BarChartExample: Story = {
  render: () => (
    <div className="h-[400px] w-[800px]">
      <ChartContainer config={chartConfig}>
        <BarChart data={data}>
          <CartesianGrid strokeDasharray="3 3" className="stroke-border" />
          <XAxis dataKey="name" stroke="#888888" fontSize={12} tickLine={false} axisLine={false} />
          <YAxis
            stroke="#888888"
            fontSize={12}
            tickLine={false}
            axisLine={false}
            tickFormatter={(value) => `$${value}`}
          />
          <ChartTooltip
            content={({ active, payload }) => {
              if (!active || !payload) return null
              return (
                <ChartTooltipContent>
                  {payload.map((item) => (
                    <div key={item.dataKey} className="flex items-center justify-between gap-2">
                      <span className="text-muted-foreground">{item.name}</span>
                      <span className="font-bold">${item.value}</span>
                    </div>
                  ))}
                </ChartTooltipContent>
              )
            }}
          />
          <Bar
            dataKey="total"
            style={% raw %}{{ fill: 'var(--color-total)', opacity: 0.8 }}{% endraw %}
          />
          <Bar
            dataKey="revenue"
            style={% raw %}{{ fill: 'var(--color-revenue)', opacity: 0.8 }}{% endraw %}
          />
        </BarChart>
      </ChartContainer>
    </div>
  ),
}

export const AreaChartExample: Story = {
  render: () => (
    <div className="h-[400px] w-[800px]">
      <ChartContainer config={chartConfig}>
        <AreaChart data={data}>
          <CartesianGrid strokeDasharray="3 3" className="stroke-border" />
          <XAxis dataKey="name" stroke="#888888" fontSize={12} tickLine={false} axisLine={false} />
          <YAxis
            stroke="#888888"
            fontSize={12}
            tickLine={false}
            axisLine={false}
            tickFormatter={(value) => `$${value}`}
          />
          <ChartTooltip
            content={({ active, payload }) => {
              if (!active || !payload) return null
              return (
                <ChartTooltipContent>
                  {payload.map((item) => (
                    <div key={item.dataKey} className="flex items-center justify-between gap-2">
                      <span className="text-muted-foreground">{item.name}</span>
                      <span className="font-bold">${item.value}</span>
                    </div>
                  ))}
                </ChartTooltipContent>
              )
            }}
          />
          <Area
            type="monotone"
            dataKey="total"
            style={% raw %}{{
              fill: 'var(--color-total)',
              opacity: 0.2,
              stroke: 'var(--color-total)',
              strokeWidth: 2,
            }}{% endraw %}
          />
          <Area
            type="monotone"
            dataKey="revenue"
            style={% raw %}{{
              fill: 'var(--color-revenue)',
              opacity: 0.2,
              stroke: 'var(--color-revenue)',
              strokeWidth: 2,
            }}{% endraw %}
          />
        </AreaChart>
      </ChartContainer>
    </div>
  ),
}
