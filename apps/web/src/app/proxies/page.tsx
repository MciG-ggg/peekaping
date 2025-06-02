import React, { useRef, useState } from "react";
import type { ProxyModel } from "@/api/types.gen";
import {
  getProxiesInfiniteOptions,
  deleteProxiesByIdMutation,
  getProxiesInfiniteQueryKey,
} from "@/api/@tanstack/react-query.gen";
import Layout from "@/layout";
import { useInfiniteQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import { Input } from "@/components/ui/input";
import { Skeleton } from "@/components/ui/skeleton";
import { Label } from "@/components/ui/label";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { toast } from "sonner";
import { Loader2, Trash } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

// ProxyCard component, similar to NotifierCard
const ProxyCard = ({
  proxy,
  onClick,
  onDelete,
}: {
  proxy: ProxyModel;
  onClick: () => void;
  onDelete?: () => void;
}) => {
  const handleDeleteClick = (e: React.MouseEvent) => {
    e.stopPropagation(); // Prevent card click when clicking delete button
    onDelete?.();
  };

  return (
    <Card
      key={proxy.id}
      className="mb-2 p-2 hover:cursor-pointer light:hover:bg-gray-100 dark:hover:bg-zinc-800"
      onClick={onClick}
    >
      <CardContent className="px-2">
        <div className="flex justify-between items-center">
          <div className="flex items-center gap-4">
            <div className="flex flex-col min-w-[100px]">
              <h3 className="font-bold mb-1">
                {proxy.host}:{proxy.port}
              </h3>
              <Badge variant={"outline"}>{proxy.protocol?.toUpperCase()} {proxy.auth ? "(auth)" : ""}</Badge>
            </div>
          </div>

          {onDelete && (
            <Button
              variant="ghost"
              size="icon"
              onClick={handleDeleteClick}
              className="text-red-500 hover:text-red-700 hover:bg-red-50 dark:hover:bg-red-950"
              aria-label={`Delete proxy ${proxy.host}`}
            >
              <Trash className="h-4 w-4" />
            </Button>
          )}
        </div>
      </CardContent>
    </Card>
  );
};

const ProxiesPage = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [search, setSearch] = useState("");
  const [showConfirmDelete, setShowConfirmDelete] = useState(false);
  const [proxyToDelete, setProxyToDelete] = useState<ProxyModel | null>(null);
  const sentinelRef = useRef<HTMLDivElement | null>(null);

  const { data, isLoading, fetchNextPage, hasNextPage, isFetchingNextPage } =
    useInfiniteQuery({
      ...getProxiesInfiniteOptions({
        query: {
          limit: 20,
          q: search || undefined,
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

  const deleteMutation = useMutation({
    ...deleteProxiesByIdMutation(),
    onSuccess: () => {
      toast.success("Proxy deleted successfully");
      queryClient.invalidateQueries({
        queryKey: getProxiesInfiniteQueryKey(),
      });
      setShowConfirmDelete(false);
      setProxyToDelete(null);
    },
    onError: () => {
      toast.error("Failed to delete proxy");
      setShowConfirmDelete(false);
      setProxyToDelete(null);
    },
  });

  const handleDeleteClick = (proxy: ProxyModel) => {
    setProxyToDelete(proxy);
    setShowConfirmDelete(true);
  };

  const handleConfirmDelete = () => {
    if (proxyToDelete?.id) {
      deleteMutation.mutate({
        path: { id: proxyToDelete.id },
      });
    }
  };

  const proxies = (data?.pages.flatMap((page) => page.data || []) || []) as ProxyModel[];

  // Infinite scroll logic
  React.useEffect(() => {
    const node = sentinelRef.current;
    if (!node) return;
    const observer = new window.IntersectionObserver(
      (entries) => {
        const [entry] = entries;
        if (entry.isIntersecting && hasNextPage && !isFetchingNextPage) {
          fetchNextPage();
        }
      },
      {
        root: null,
        rootMargin: "0px",
        threshold: 1.0,
      }
    );
    observer.observe(node);
    return () => {
      observer.unobserve(node);
    };
  }, [fetchNextPage, hasNextPage, isFetchingNextPage]);

  return (
    <Layout pageName="Proxies" onCreate={() => navigate("/proxies/new")}>
      <div>
        <div className="mb-4 flex justify-end gap-4">
          <div className="flex flex-col gap-1">
            <Label htmlFor="search-proxies">Search</Label>
            <Input
              id="search-proxies"
              placeholder="Search proxies by host..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-[400px]"
            />
          </div>
        </div>
        {proxies.length === 0 && isLoading && (
          <div className="flex flex-col space-y-2 mb-2">
            {Array.from({ length: 7 }, (_, id) => (
              <Skeleton className="h-[68px] w-full rounded-xl" key={id} />
            ))}
          </div>
        )}
        {/* Proxies list */}
        {proxies.map((proxy) => (
          <ProxyCard
            key={proxy.id}
            proxy={proxy}
            onClick={() => navigate(`/proxies/edit/${proxy.id}`)}
            onDelete={() => handleDeleteClick(proxy)}
          />
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

      <AlertDialog open={showConfirmDelete} onOpenChange={setShowConfirmDelete}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
            <AlertDialogDescription>
              This action cannot be undone. This will permanently delete the proxy {proxyToDelete?.host}:{proxyToDelete?.port}.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel onClick={() => setShowConfirmDelete(false)}>
              Cancel
            </AlertDialogCancel>
            <AlertDialogAction
              onClick={handleConfirmDelete}
              disabled={deleteMutation.isPending}
              className="bg-red-600 hover:bg-red-700 focus:ring-red-600"
            >
              {deleteMutation.isPending && <Loader2 className="animate-spin mr-2 h-4 w-4" />}
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </Layout>
  );
};

export default ProxiesPage;
