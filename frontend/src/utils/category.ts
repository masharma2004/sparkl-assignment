const DEFAULT_CATEGORY = 'General'

export function normalizeCategory(category: string | null | undefined) {
  const value = typeof category === 'string' ? category.trim() : ''
  return value.length > 0 ? value : DEFAULT_CATEGORY
}

export function matchesCategoryFilter(
  category: string | null | undefined,
  filter: string,
) {
  if (filter === 'all') {
    return true
  }

  return normalizeCategory(category) === filter
}
