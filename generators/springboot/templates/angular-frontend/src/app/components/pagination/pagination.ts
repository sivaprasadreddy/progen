import { Component, ChangeDetectionStrategy, input, output, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { getPageNumbers } from '../../utils/pagination.util';

/**
 * Display variant for the pagination component.
 * - 'default': Simple pagination without background wrapper
 * - 'card': Pagination wrapped in a white card with shadow
 */
export type PaginationVariant = 'default' | 'card';

/**
 * Reusable pagination component for navigating through pages.
 * Supports previous/next buttons, page number buttons with ellipsis, and page info display.
 * All page numbers are 1-indexed (user-facing).
 */
@Component({
  selector: 'app-pagination',
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [CommonModule],
  templateUrl: './pagination.html',
})
export class Pagination {
  /**
   * Current page number (1-indexed).
   */
  readonly currentPage = input.required<number>();

  /**
   * Total number of pages.
   */
  readonly totalPages = input.required<number>();

  /**
   * Whether to show the page info text (e.g., "Page 1 of 10").
   * Default: true
   */
  readonly showPageInfo = input<boolean>(true);

  /**
   * Visual variant of the pagination component.
   * Default: 'default'
   */
  readonly variant = input<PaginationVariant>('default');

  /**
   * Event emitted when the user navigates to a different page.
   * Emits the new page number (1-indexed).
   */
  readonly pageChange = output<number>();

  /**
   * Computed: Array of page numbers to display, with -1 representing ellipsis.
   */
  readonly pageNumbers = computed(() => {
    return getPageNumbers(this.currentPage(), this.totalPages());
  });

  /**
   * Computed: Whether the Previous button should be disabled.
   */
  readonly isPreviousDisabled = computed(() => {
    return this.currentPage() === 1;
  });

  /**
   * Computed: Whether the Next button should be disabled.
   */
  readonly isNextDisabled = computed(() => {
    return this.currentPage() === this.totalPages();
  });

  /**
   * Handles navigation to the previous page.
   */
  protected handlePreviousPage(): void {
    if (!this.isPreviousDisabled()) {
      this.pageChange.emit(this.currentPage() - 1);
    }
  }

  /**
   * Handles navigation to the next page.
   */
  protected handleNextPage(): void {
    if (!this.isNextDisabled()) {
      this.pageChange.emit(this.currentPage() + 1);
    }
  }

  /**
   * Handles navigation to a specific page.
   */
  protected handleGoToPage(page: number): void {
    if (page >= 1 && page <= this.totalPages() && page !== this.currentPage()) {
      this.pageChange.emit(page);
    }
  }
}
