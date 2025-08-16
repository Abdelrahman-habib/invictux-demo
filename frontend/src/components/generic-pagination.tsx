"use client";

import { Loader2 } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";

interface GenericPaginationProps {
  canGoNext: boolean;
  canGoPrev: boolean;
  onGoToNextPage: () => void;
  onGoToPrevPage: () => void;
  onGoToPage: (page: number) => void;
  currentPage: number;
  totalPages: number;
  isLoading: boolean;
  maxPagesToShow?: number;
}

export function GenericPagination({
  canGoNext,
  canGoPrev,
  onGoToNextPage,
  onGoToPrevPage,
  onGoToPage,
  currentPage,
  totalPages,
  isLoading,
  maxPagesToShow = 3,
}: GenericPaginationProps) {
  // Calculate the range of pages to show
  const halfPages = Math.floor(maxPagesToShow / 2);
  let startPage = Math.max(1, currentPage - halfPages);
  const endPage = Math.min(totalPages, startPage + maxPagesToShow - 1);

  // Adjust start page if we're near the end
  if (endPage - startPage + 1 < maxPagesToShow) {
    startPage = Math.max(1, endPage - maxPagesToShow + 1);
  }

  const pagesToShow: (number | "ellipsis")[] = [];

  // Always show first page if we're not starting from it
  if (startPage > 1) {
    pagesToShow.push(1);
    if (startPage > 2) {
      pagesToShow.push("ellipsis");
    }
  }

  // Add the main range of pages
  for (let i = startPage; i <= endPage; i++) {
    pagesToShow.push(i);
  }

  // Always show last page if we're not ending with it
  if (endPage < totalPages) {
    if (endPage < totalPages - 1) {
      pagesToShow.push("ellipsis");
    }
    pagesToShow.push(totalPages);
  }

  return (
    <Pagination>
      <PaginationContent>
        <PaginationItem>
          <PaginationPrevious
            onClick={() => !(isLoading || !canGoPrev) && onGoToPrevPage()}
            className={
              isLoading || !canGoPrev ? "pointer-events-none opacity-50" : ""
            }
            href={isLoading || !canGoPrev ? undefined : "#"} // Prevent navigation when disabled
          />
        </PaginationItem>

        {pagesToShow.map((page, index) =>
          page === "ellipsis" ? (
            <PaginationItem key={`ellipsis-${index}`}>
              <PaginationEllipsis />
            </PaginationItem>
          ) : (
            <PaginationItem key={page}>
              <PaginationLink
                onClick={() =>
                  !(isLoading || page === currentPage) &&
                  onGoToPage(page as number)
                }
                isActive={page === currentPage}
                className={
                  isLoading || page === currentPage
                    ? "pointer-events-none opacity-50"
                    : ""
                }
                href={isLoading || page === currentPage ? undefined : "#"} // Prevent navigation when disabled
              >
                {page}
              </PaginationLink>
            </PaginationItem>
          )
        )}

        <PaginationItem>
          <PaginationNext
            onClick={() => !(isLoading || !canGoNext) && onGoToNextPage()}
            className={
              isLoading || !canGoNext ? "pointer-events-none opacity-50" : ""
            }
            href={isLoading || !canGoNext ? undefined : "#"} // Prevent navigation when disabled
          >
            {isLoading ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <span>Next</span>
            )}
          </PaginationNext>
        </PaginationItem>
      </PaginationContent>
    </Pagination>
  );
}

export function GenericPaginationSkeleton() {
  return (
    <div className="flex items-center justify-center gap-2">
      <Skeleton className="h-9 w-9" />
      <Skeleton className="h-9 w-9" />
      <Skeleton className="h-9 w-9" />
      <Skeleton className="h-9 w-9" />
    </div>
  );
}
