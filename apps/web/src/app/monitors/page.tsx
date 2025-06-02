import { useInfiniteQuery, useQueryClient } from "@tanstack/react-query";
import {
  getMonitorsByIdHeartbeatsQueryKey,
  getMonitorsInfiniteOptions,
} from "@/api/@tanstack/react-query.gen";
import Layout from "@/layout";
import { useNavigate } from "react-router-dom";
import type {
  HeartbeatModel,
  MonitorModel,
  UtilsApiResponseArrayHeartbeatModel,
} from "@/api";
import { useWebSocket, WebSocketStatus } from "@/context/WebsocketContext";
import { useEffect, useState, useRef, useCallback } from "react";
import { useDebounce } from "@/hooks/useDebounce";
import { Input } from "@/components/ui/input";
import { Skeleton } from "@/components/ui/skeleton";
import MonitorCard from "./components/monitor-card";
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from "@/components/ui/select";
import { Label } from "@/components/ui/label";

const MonitorsPage = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  // Add state for search query
  const [search, setSearch] = useState("");
  const debouncedSearch = useDebounce(search, 400);

  // Add state for active filter
  const [activeFilter, setActiveFilter] = useState<
    "all" | "active" | "inactive"
  >("all");
  const [statusFilter, setStatusFilter] = useState<
    "all" | "up" | "down" | "maintenance"
  >("all");

  const { data, isLoading, fetchNextPage, hasNextPage, isFetchingNextPage } =
    useInfiniteQuery({
      ...getMonitorsInfiniteOptions({
        query: {
          limit: 20,
          q: debouncedSearch || undefined,
          active:
            activeFilter === "all"
              ? undefined
              : activeFilter === "active"
              ? true
              : false,
          status:
            statusFilter === "all"
              ? undefined
              : statusFilter === "up"
              ? 1
              : statusFilter === "down"
              ? 0
              : statusFilter === "maintenance"
              ? 3
              : undefined,
        },
      }),
      getNextPageParam: (lastPage, pages) => {
        const lastLength = lastPage.data?.length || 0;
        if (lastLength < 20) return undefined;
        return pages.length;
      },
      initialPageParam: 0,
      enabled: true,
    });

  const monitors = (data?.pages.flatMap((page) => page.data || []) ||
    []) as MonitorModel[];

  const { socket, status: socketStatus } = useWebSocket();
  const subscribedRef = useRef(false);

  useEffect(() => {
    if (!socket || socketStatus !== WebSocketStatus.CONNECTED) return;
    if (subscribedRef.current) return;
    subscribedRef.current = true;

    const roomName = "monitor:all";

    const handleHeartbeat = (newHeartbeat: HeartbeatModel) => {
      queryClient.setQueryData(
        getMonitorsByIdHeartbeatsQueryKey({
          path: {
            id: newHeartbeat.monitor_id!,
          },
          query: {
            limit: 50,
            reverse: true,
          },
        }),
        (oldData: UtilsApiResponseArrayHeartbeatModel) => {
          if (!oldData) return oldData;
          return {
            ...oldData,
            data: [...(oldData.data || []), newHeartbeat].slice(-50),
          };
        }
      );
    };

    socket.on(`${roomName}:heartbeat`, handleHeartbeat);
    socket.emit("join_room", roomName);
    console.log("Subscribed to heartbeat", roomName);

    // return () => {
    //   socket.off("heartbeat", handleHeartbeat);
    //   if (socketStatus === WebSocketStatus.CONNECTED) {
    //     socket.emit("leave_room", roomName);
    //   }
    // };
  }, [socket, socketStatus, queryClient]);

  // Infinite scroll logic
  const sentinelRef = useRef<HTMLDivElement | null>(null);

  const handleObserver = useCallback(
    (entries: IntersectionObserverEntry[]) => {
      const [entry] = entries;
      if (entry.isIntersecting && hasNextPage && !isFetchingNextPage) {
        fetchNextPage();
      }
    },
    [fetchNextPage, hasNextPage, isFetchingNextPage]
  );

  useEffect(() => {
    const node = sentinelRef.current;
    if (!node) return;
    const observer = new window.IntersectionObserver(handleObserver, {
      root: null,
      rootMargin: "0px",
      threshold: 1.0,
    });
    observer.observe(node);
    return () => {
      observer.unobserve(node);
    };
  }, [handleObserver]);

  return (
    <Layout
      pageName="Monitors"
      onCreate={() => {
        navigate("/monitors/new");
      }}
    >
      <div>
        <div className="mb-4 flex justify-end gap-4">
          <div className="flex flex-col gap-1">
            <Label htmlFor="search-monitors">Search</Label>
            <Input
              id="search-monitors"
              placeholder="Search monitors by name..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-[400px]"
            />
          </div>
          <div className="flex flex-col gap-1">
            <Label htmlFor="active-filter">Active</Label>
            <Select
              value={activeFilter}
              onValueChange={(v) =>
                setActiveFilter(v as "all" | "active" | "inactive")
              }
            >
              <SelectTrigger className="w-[140px]">
                <SelectValue placeholder="Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All</SelectItem>
                <SelectItem value="active">Active</SelectItem>
                <SelectItem value="inactive">Inactive</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div className="flex flex-col gap-1">
            <Label htmlFor="status-filter">Monitor Status</Label>
            <Select
              value={statusFilter}
              onValueChange={(v) =>
                setStatusFilter(v as "all" | "up" | "down" | "maintenance")
              }
            >
              <SelectTrigger className="w-[160px]">
                <SelectValue placeholder="Monitor Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All</SelectItem>
                <SelectItem value="up">Up</SelectItem>
                <SelectItem value="down">Down</SelectItem>
                <SelectItem value="maintenance">Maintenance</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>

        {monitors.length === 0 && isLoading && (
          <div className="flex flex-col space-y-2 mb-2">
            {Array.from({ length: 7 }, (_, id) => (
              <Skeleton className="h-[68px] w-full rounded-xl" key={id} />
            ))}
          </div>
        )}

        {/* Monitors list */}
        {monitors.map((monitor) => (
          <MonitorCard key={monitor.id} monitor={monitor} />
        ))}
        {/* Sentinel for infinite scroll */}
        <div ref={sentinelRef} style={{ height: 1 }} />
        {isFetchingNextPage && (
          <div className="flex flex-col space-y-2 mb-2">
            {Array.from({ length: 3 }, (_, i) => (
              <Skeleton key={i} className="h-[68px] w-full rounded-xl" />
            ))}
          </div>
        )}
      </div>
    </Layout>
  );
};

export default MonitorsPage;
