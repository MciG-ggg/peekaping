"use client";

import * as React from "react";
import {
  CartesianGrid,
  Line,
  LineChart,
  ReferenceArea,
  XAxis,
  YAxis,
} from "recharts";

import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  type ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useQuery } from "@tanstack/react-query";
import { getMonitorsByIdChartpointsOptions } from "@/api/@tanstack/react-query.gen";
import type { HeartbeatChartPoint } from "@/api";
import { useTimezone } from '../context/TimezoneContext';
import { formatDateToTimezone } from '../lib/formatDateToTimezone';
// import { useMemo } from "react";

function getStatusRanges(data: HeartbeatChartPoint[]) {
  if (!data.length) return [];
  const ranges = [];
  let startIdx = 0;
  let currentColor = getColor(data[0]);

  function getColor(point: HeartbeatChartPoint) {
    if (point.up !== 0 && point.down === 0) return "none";
    if (point.up !== 0 && point.down !== 0) return "rgba(255, 215, 0, 0.15)";
    if (point.up === 0 && point.down !== 0) return "rgba(255, 0, 0, 0.15)";
    return null;
  }

  for (let i = 1; i < data.length; i++) {
    const color = getColor(data[i]);
    if (color !== currentColor) {
      if (currentColor && data[startIdx].timestamp !== data[i].timestamp) {
        ranges.push({
          x1: data[startIdx].timestamp,
          x2: data[i].timestamp,
          color: currentColor,
        });
      }
      startIdx = i;
      currentColor = color;
    }
  }
  // Push the last range, filter zero-width
  if (
    currentColor &&
    data[startIdx].timestamp !== data[data.length - 1].timestamp
  ) {
    ranges.push({
      x1: data[startIdx].timestamp,
      x2: data[data.length - 1].timestamp,
      color: currentColor,
    });
  }
  return ranges;
}

const chartConfig = {
  visitors: {
    label: "Visitors",
  },
} satisfies ChartConfig;

export function Chart({ id }: { id: string }) {
  const [timeRange, setTimeRange] = React.useState<
    "30m" | "3h" | "6h" | "24h" | "1week"
  >("30m");
  const { timezone } = useTimezone();

  // const filteredData = chartData.filter((item) => {
  //   const date = new Date(item.date);
  //   const referenceDate = new Date("2024-06-30");
  //   let daysToSubtract = 90;
  //   if (timeRange === "30d") {
  //     daysToSubtract = 30;
  //   } else if (timeRange === "7d") {
  //     daysToSubtract = 7;
  //   }
  //   const startDate = new Date(referenceDate);
  //   startDate.setDate(startDate.getDate() - daysToSubtract);
  //   return date >= startDate;
  // });

  const {
    data: chartDataRaw,
    // error: chartDataError,
    // isLoading: chartDataIsLoading,
  } = useQuery({
    ...getMonitorsByIdChartpointsOptions({
      path: {
        id: id!,
      },
      query: {
        period: timeRange,
      },
    }),
    refetchInterval: 5 * 1000,
    enabled: !!id,
  });

  const chartData = chartDataRaw?.data || [];

  const statusRanges = getStatusRanges(chartData);

  // stats: min, max, avg, median
  const stats = React.useMemo(() => {
    if (!chartData.length) return { min: 0, max: 0, avg: 0, median: 0 };
    const worthPoints = chartData.filter((e) => !(e.up === 0 && e.down === 0));
    // Filter only points where up is not zero and avgPing is a number
    const max = Math.max(...worthPoints.map((el) => el.maxPing!));
    const min = Math.min(...worthPoints.map((el) => el.minPing!));
    const avg = Math.max(...worthPoints.map((el) => el.avgPing!));

    return { min, max, avg };
  }, [chartData]);

  // Create an array with labels for displaying stats
  const statsArray = [
    { key: "min", label: "Minimum", value: stats.min || 0 },
    { key: "max", label: "Maximum", value: stats.max || 0 },
    { key: "avg", label: "Average", value: stats.avg || 0 },
  ];

  return (
    <Card>
      <CardHeader className="flex items-center gap-2 space-y-0 border-b sm:flex-row">
        <div className="grid flex-1 gap-1 text-center sm:text-left">
          <CardTitle>Response time</CardTitle>
          {/* <CardDescription>
            Showing total visitors for the last 3 months
          </CardDescription> */}
        </div>

        <Select value={timeRange} onValueChange={setTimeRange as any}>
          <SelectTrigger
            className="w-[160px] rounded-lg sm:ml-auto"
            aria-label="Select a value"
          >
            <SelectValue placeholder="Last 3 months" />
          </SelectTrigger>

          <SelectContent className="rounded-xl">
            <SelectItem value="30m" className="rounded-lg">
              Last 30 minutes
            </SelectItem>
            <SelectItem value="3h" className="rounded-lg">
              Last 3 hours
            </SelectItem>
            <SelectItem value="6h" className="rounded-lg">
              Last 6 hours
            </SelectItem>
            <SelectItem value="24h" className="rounded-lg">
              Last 24 hours
            </SelectItem>
            <SelectItem value="1week" className="rounded-lg">
              Last 7 days
            </SelectItem>
          </SelectContent>
        </Select>
      </CardHeader>

      <CardContent className="px-2 pt-2 sm:px-6 sm:pt-6">
        <ChartContainer
          config={chartConfig}
          className="aspect-auto h-[250px] w-full"
        >
          <LineChart
            data={chartData.map((el) => ({
              ...el,
              avgPing: el.up ? el.avgPing : null,
              minPing: el.up ? el.minPing : null,
              maxPing: el.up ? el.maxPing : null,
            }))}
            accessibilityLayer
            margin={{
              left: 12,
              right: 12,
            }}
          >
            <CartesianGrid vertical={false} />
            <XAxis
              dataKey="timestamp"
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              // minTickGap={32}
              tickFormatter={(timestamp) => {
                return formatDateToTimezone(timestamp, timezone, {
                  month: 'short',
                  day: 'numeric',
                  hour: '2-digit',
                  minute: '2-digit',
                });
              }}
            />

            <YAxis
              tickLine={false}
              axisLine={false}
              label={{
                value: "Resp. Time (ms)", // the text you want to show
                angle: -90, // -90 = reading bottom-to-top
                position: "insideLeft", // or 'outsideLeft', 'insideRight', …
                offset: 0, // fine-tune distance from the axis
                style: { textAnchor: "middle" }, // keep it centred on the axis
              }}
            />

            <ChartTooltip
              cursor={false}
              content={
                <ChartTooltipContent
                  labelFormatter={(value) => {
                    return formatDateToTimezone(value, timezone, {
                      month: 'short',
                      day: 'numeric',
                      hour: '2-digit',
                      minute: '2-digit',
                    });
                  }}
                  indicator="dot"
                />
              }
            />

            <Line
              dataKey="avgPing"
              type="monotone"
              stroke="var(--chart-4)"
              strokeWidth={2}
              dot={false}
              connectNulls={false}
            />

            <Line
              dataKey="minPing"
              type="monotone"
              stroke="var(--chart-3)"
              strokeWidth={2}
              dot={false}
              connectNulls={false}
            />

            <Line
              dataKey="maxPing"
              type="monotone"
              stroke="var(--chart-2)"
              strokeWidth={2}
              dot={false}
              connectNulls={false}
            />

            {statusRanges
              .filter((e) => e.color !== "none")
              .map(({ x1, x2, color }, idx) => (
                <ReferenceArea
                  key={idx}
                  x1={x1}
                  x2={x2}
                  strokeOpacity={0}
                  fill={color}
                />
              ))}

            {/* <ChartLegend content={<ChartLegendContent />} /> */}
          </LineChart>
        </ChartContainer>
      </CardContent>

      <CardFooter className="flex flex-col items-stretch space-y-0 border-t p-0">
        <div className="grid grid-cols-1 md:grid-cols-3">
          {statsArray.map((item) => {
            return (
              <div
                key={item.key}
                className="flex flex-1 flex-col justify-center gap-1  px-6 py-4 text-left even:border-l sm:border-l sm:border-t-0 sm:px-8 sm:py-6"
              >
                <span className="text-xs text-muted-foreground">
                  {item.label}
                </span>
                <span className="text-lg font-bold leading-none sm:text-3xl">
                  {item.value.toLocaleString()}{" "}
                  <span className="text-sm font-normal text-muted-foreground">
                    ms
                  </span>
                </span>
              </div>
            );
          })}
        </div>
      </CardFooter>
    </Card>
  );
}
