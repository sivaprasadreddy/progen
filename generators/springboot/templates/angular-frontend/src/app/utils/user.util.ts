/**
 * Generates user initials from a full name.
 * - If the full name has 2 or more words, returns first letter of first and last word
 * - If the full name has 1 word, returns first letter of that word
 * - If no full name is provided but a fallback is given, returns first letter of fallback
 * - Returns empty string if no name is available
 *
 * @param fullName - The user's full name (e.g., "John Doe")
 * @param fallback - Optional fallback string (e.g., username) to use if fullName is empty
 * @returns User initials in uppercase (e.g., "JD")
 *
 * @example
 * getUserInitials('John Doe') // returns 'JD'
 * getUserInitials('John') // returns 'J'
 * getUserInitials('', 'johndoe') // returns 'J'
 * getUserInitials('') // returns ''
 */
export function getUserInitials(fullName: string, fallback?: string): string {
  // Try fullName first
  if (fullName && fullName.trim()) {
    const names = fullName.trim().split(/\s+/); // Split by any whitespace

    if (names.length >= 2) {
      // Return first letter of first and last name
      const firstInitial = names[0].charAt(0);
      const lastInitial = names[names.length - 1].charAt(0);
      return `${firstInitial}${lastInitial}`.toUpperCase();
    }

    // Return first letter of single name
    return names[0].charAt(0).toUpperCase();
  }

  // Fall back to fallback string if provided
  if (fallback && fallback.trim()) {
    return fallback.trim().charAt(0).toUpperCase();
  }

  // No name available
  return '';
}
