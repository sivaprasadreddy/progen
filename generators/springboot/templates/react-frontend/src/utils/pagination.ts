/**
 * Pagination utility functions for calculating page numbers and navigation.
 * All page numbers are 1-indexed (user-facing).
 */

/**
 * Calculates an array of page numbers to display in pagination controls.
 * Returns -1 for ellipsis (...) positions.
 *
 * Logic:
 * - If total pages <= 7: Show all pages
 * - If current page is near start (< 5): Show first 5 pages, ellipsis, last page
 * - If current page is near end (> total - 4): Show first page, ellipsis, last 5 pages
 * - Otherwise: Show first page, ellipsis, current +/- 1 pages, ellipsis, last page
 *
 * @param currentPage - One-based current page number (1, 2, 3, ...)
 * @param totalPages - Total number of pages
 * @returns Array of page numbers (one-based), with -1 representing ellipsis
 *
 * @example
 * getPageNumbers(1, 3) // [1, 2, 3]
 * getPageNumbers(3, 10) // [1, -1, 2, 3, 4, -1, 10]
 * getPageNumbers(9, 10) // [1, -1, 6, 7, 8, 9, 10]
 */
export function getPageNumbers(currentPage: number, totalPages: number): number[] {
  const pages: number[] = [];

  // Show all pages if total is 7 or less
  if (totalPages <= 7) {
    for (let i = 1; i <= totalPages; i++) {
      pages.push(i);
    }
    return pages;
  }

  // Current page is near the start
  if (currentPage < 5) {
    for (let i = 1; i <= 5; i++) {
      pages.push(i);
    }
    pages.push(-1); // Ellipsis
    pages.push(totalPages); // Last page
    return pages;
  }

  // Current page is near the end
  if (currentPage > totalPages - 4) {
    pages.push(1); // First page
    pages.push(-1); // Ellipsis
    for (let i = totalPages - 4; i <= totalPages; i++) {
      pages.push(i);
    }
    return pages;
  }

  // Current page is in the middle
  pages.push(1); // First page
  pages.push(-1); // Ellipsis
  for (let i = currentPage - 1; i <= currentPage + 1; i++) {
    pages.push(i);
  }
  pages.push(-1); // Ellipsis
  pages.push(totalPages); // Last page

  return pages;
}

/**
 * Checks if navigation to the previous page is possible.
 *
 * @param currentPage - One-based current page number
 * @returns True if previous page navigation is possible
 */
export function canNavigateToPrevious(currentPage: number): boolean {
  return currentPage > 1;
}

/**
 * Checks if navigation to the next page is possible.
 *
 * @param currentPage - One-based current page number
 * @param totalPages - Total number of pages
 * @returns True if next page navigation is possible
 */
export function canNavigateToNext(currentPage: number, totalPages: number): boolean {
  return currentPage < totalPages;
}

/**
 * Calculates the previous page number.
 *
 * @param currentPage - One-based current page number
 * @returns Previous page number, or current page if at start
 */
export function getPreviousPage(currentPage: number): number {
  return Math.max(1, currentPage - 1);
}

/**
 * Calculates the next page number.
 *
 * @param currentPage - One-based current page number
 * @param totalPages - Total number of pages
 * @returns Next page number, or current page if at end
 */
export function getNextPage(currentPage: number, totalPages: number): number {
  return Math.min(totalPages, currentPage + 1);
}

/**
 * Validates and clamps a page number to valid range.
 *
 * @param page - Page number to validate
 * @param totalPages - Total number of pages
 * @returns Clamped page number within valid range [1, totalPages]
 */
export function clampPageNumber(page: number, totalPages: number): number {
  return Math.max(1, Math.min(totalPages, page));
}
