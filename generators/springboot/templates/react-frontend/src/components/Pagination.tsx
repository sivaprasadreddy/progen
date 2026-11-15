import { getPageNumbers } from '../utils/pagination';
import { cn } from '../utils/classNames';

/**
 * Display variant for the pagination component.
 * - 'default': Simple pagination without background wrapper
 * - 'card': Pagination wrapped in a white card with shadow
 */
export type PaginationVariant = 'default' | 'card';

interface PaginationProps {
  /** Current page number (1-indexed). */
  currentPage: number;
  /** Total number of pages. */
  totalPages: number;
  /** Whether to show the page info text (e.g., "Page 1 of 10"). Default: true */
  showPageInfo?: boolean;
  /** Visual variant of the pagination component. Default: 'default' */
  variant?: PaginationVariant;
  /** Called with the new page number (1-indexed) when the user navigates. */
  onPageChange: (page: number) => void;
}

/**
 * Reusable pagination component for navigating through pages.
 * Supports previous/next buttons, page number buttons with ellipsis, and page info display.
 */
export default function Pagination({
  currentPage,
  totalPages,
  showPageInfo = true,
  variant = 'default',
  onPageChange,
}: PaginationProps) {
  if (totalPages <= 1) {
    return null;
  }

  const isPreviousDisabled = currentPage === 1;
  const isNextDisabled = currentPage === totalPages;

  const goToPage = (page: number) => {
    if (page >= 1 && page <= totalPages && page !== currentPage) {
      onPageChange(page);
    }
  };

  return (
    <div className={cn(variant === 'card' && 'bg-white rounded-lg shadow-md p-4')}>
      <div className="flex justify-center items-center gap-2">
        <button
          type="button"
          onClick={() => goToPage(currentPage - 1)}
          disabled={isPreviousDisabled}
          className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          aria-label="Previous page"
        >
          Previous
        </button>

        <div className="flex gap-1">
          {getPageNumbers(currentPage, totalPages).map((pageNum, index) =>
            pageNum === -1 ? (
              <span key={`ellipsis-${index}`} className="px-3 py-2 text-gray-500">
                ...
              </span>
            ) : (
              <button
                type="button"
                key={pageNum}
                onClick={() => goToPage(pageNum)}
                className={cn(
                  'px-4 py-2 border border-gray-300 rounded-lg transition-colors',
                  pageNum === currentPage ? 'bg-blue-600 text-white' : 'hover:bg-gray-50',
                )}
                aria-label={`Go to page ${pageNum}`}
                aria-current={pageNum === currentPage ? 'page' : undefined}
              >
                {pageNum}
              </button>
            ),
          )}
        </div>

        <button
          type="button"
          onClick={() => goToPage(currentPage + 1)}
          disabled={isNextDisabled}
          className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          aria-label="Next page"
        >
          Next
        </button>
      </div>

      {showPageInfo && (
        <div className="text-center mt-4 text-sm text-gray-600">
          Page {currentPage} of {totalPages}
        </div>
      )}
    </div>
  );
}
