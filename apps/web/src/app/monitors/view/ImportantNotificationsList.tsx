import { useInfiniteQuery } from "@tanstack/react-query";
import { getMonitorsByIdHeartbeatsInfiniteOptions } from "@/api/@tanstack/react-query.gen";
import { useRef, useEffect } from "react";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import { TypographyH4 } from "@/components/ui/typography";
import { formatDateToTimezone } from '../../../lib/formatDateToTimezone';
import { useTimezone } from '../../../context/TimezoneContext';

const ImportantNotificationsList = ({ monitorId }: { monitorId: string }) => {
  const sentinelRef = useRef<HTMLDivElement | null>(null);
  const { timezone } = useTimezone();
  const {
    data,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
  } = useInfiniteQuery({
    ...getMonitorsByIdHeartbeatsInfiniteOptions({
      path: { id: monitorId },
      query: { important: true, limit: 20 },
    }),
    getNextPageParam: (lastPage, pages) => {
      if ((lastPage.data?.length ?? 0) < 20) return undefined;
      return pages.length;
    },
    initialPageParam: 0,
    enabled: !!monitorId,
    staleTime: 0,
  });

  useEffect(() => {
    const node = sentinelRef.current;
    if (!node) return;
    const observer = new window.IntersectionObserver(
      (entries) => {
        const [entry] = entries;
        if (entry.isIntersecting && hasNextPage && !isFetchingNextPage) {
          fetchNextPage();
        }
      },
      { root: null, rootMargin: "0px", threshold: 1.0 }
    );
    observer.observe(node);
    return () => observer.unobserve(node);
  }, [fetchNextPage, hasNextPage, isFetchingNextPage]);

  const importantHeartbeats = data?.pages.flatMap((page) => page.data || []) ?? [];

  return (
    <div className="mb-6 mt-6">
      <TypographyH4 className="mb-2">Important Notifications</TypographyH4>
      {importantHeartbeats.length === 0 && isLoading && <div>Loading...</div>}
      {importantHeartbeats.length === 0 && !isLoading && (
        <div className="text-muted-foreground">No important notifications found.</div>
      )}

      {importantHeartbeats.map((hb) => (
        <Card key={hb.id} className="mb-2 p-2">
          <CardContent className="p-2 flex flex-col gap-2">
            <div className="flex items-center gap-2">
              <Badge
                variant={hb.status === 1 ? "default" : "destructive"}
                className={
                  hb.status === 1
                    ? "bg-green-500 text-white border-green-600"
                    : "bg-red-500 text-white border-red-600"
                }
              >
                {hb.status === 1 ? "Up" : "Down"}
              </Badge>
              <span className="text-xs text-muted-foreground">
                {hb.time && formatDateToTimezone(hb.time, timezone)}
              </span>
            </div>
            <div className="font-medium text-sm">{hb.msg}</div>
            <div className="flex flex-wrap gap-4 text-xs text-muted-foreground">
              <span>Ping: <span className="text-foreground">{hb.ping} ms</span></span>
              <span>Retries: <span className="text-foreground">{hb.retries}</span></span>
              {typeof hb.down_count !== "undefined" && (
                <span>Down count: <span className="text-foreground">{hb.down_count}</span></span>
              )}
              <span>Notified: <span className="text-foreground">{hb.notified ? "Yes" : "No"}</span></span>
            </div>
          </CardContent>
        </Card>
      ))}
      <div ref={sentinelRef} style={{ height: 1 }} />
      {isFetchingNextPage && <div>Loading more...</div>}
    </div>
  );
};

export default ImportantNotificationsList;
