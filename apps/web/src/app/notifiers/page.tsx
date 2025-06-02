import React from "react";
import type { NotificationModel } from "@/api";
import {
  getNotificationsInfiniteOptions,
  deleteNotificationsByIdMutation,
  getNotificationsInfiniteQueryKey
} from "@/api/@tanstack/react-query.gen";
import Layout from "@/layout";
import { useInfiniteQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import { useState, useRef } from "react";
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
import NotifierCard from "./notifier-card";
import { useDebounce } from "@/hooks/useDebounce";
import { toast } from "sonner";
import { Loader2 } from "lucide-react";

const NotifiersPage = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [search, setSearch] = useState("");
  const debouncedSearch = useDebounce(search, 400);
  const sentinelRef = useRef<HTMLDivElement | null>(null);
  const [showConfirmDelete, setShowConfirmDelete] = useState(false);
  const [notifierToDelete, setNotifierToDelete] = useState<NotificationModel | null>(null);

  const { data, isLoading, fetchNextPage, hasNextPage, isFetchingNextPage } =
    useInfiniteQuery({
      ...getNotificationsInfiniteOptions({
        query: {
          limit: 20,
          q: debouncedSearch || undefined,
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
    ...deleteNotificationsByIdMutation(),
    onSuccess: () => {
      toast.success("Notifier deleted successfully");
      // Invalidate and refetch the notifiers list
      queryClient.invalidateQueries({
        queryKey: getNotificationsInfiniteQueryKey(),
      });
      setShowConfirmDelete(false);
      setNotifierToDelete(null);
    },
    onError: () => {
      toast.error("Failed to delete notifier");
      setShowConfirmDelete(false);
      setNotifierToDelete(null);
    },
  });

  const handleDeleteClick = (notifier: NotificationModel) => {
    setNotifierToDelete(notifier);
    setShowConfirmDelete(true);
  };

  const handleConfirmDelete = () => {
    if (notifierToDelete?.id) {
      deleteMutation.mutate({
        path: { id: notifierToDelete.id },
      });
    }
  };

  const notifications = (data?.pages.flatMap((page) => page.data || []) ||
    []) as NotificationModel[];

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
    <Layout pageName="Notifiers" onCreate={() => navigate("/notifiers/new")}>
      <div>
        <div className="mb-4 flex justify-end gap-4">
          <div className="flex flex-col gap-1">
            <Label htmlFor="search-notifiers">Search</Label>
            <Input
              id="search-notifiers"
              placeholder="Search notifiers by name..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-[400px]"
            />
          </div>
          {/* Add more filters here if needed */}
        </div>
        {notifications.length === 0 && isLoading && (
          <div className="flex flex-col space-y-2 mb-2">
            {Array.from({ length: 7 }, (_, id) => (
              <Skeleton className="h-[68px] w-full rounded-xl" key={id} />
            ))}
          </div>
        )}
        {/* Notifiers list */}
        {notifications.map((notifier) => (
          <NotifierCard
            key={notifier.id}
            notifier={notifier}
            onClick={() => navigate(`/notifiers/${notifier.id}/edit`)}
            onDelete={() => handleDeleteClick(notifier)}
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
              This action cannot be undone. This will permanently delete the notifier "{notifierToDelete?.name}" and remove it from all monitors that use it.
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

export default NotifiersPage;
