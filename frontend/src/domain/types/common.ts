import { z } from "zod";

// Pagination state schema
export const paginationStateSchema = z.object({
  currentPage: z.number().min(1).default(1),
  limit: z.number().min(1).max(100).default(20),
  searchQuery: z.string().default(""),
});

// Pagination result schema
export const paginationResultSchema = z.object({
  totalCount: z.number(),
  currentPage: z.number(),
  totalPages: z.number(),
  hasNextPage: z.boolean(),
  hasPrevPage: z.boolean(),
});

// Device filters schema
export const deviceFiltersSchema = z.object({
  deviceType: z.string().optional(),
  vendor: z.string().optional(),
  status: z.string().optional(),
  tags: z.array(z.string()).optional(),
});

// Security filters schema
export const securityFiltersSchema = z.object({
  severity: z.string().optional(),
  status: z.string().optional(),
  deviceId: z.string().optional(),
  checkType: z.string().optional(),
  dateRange: z
    .object({
      startDate: z.date(),
      endDate: z.date(),
    })
    .optional(),
});

// Sort options schema
export const sortOptionsSchema = z.object({
  field: z.string(),
  direction: z.enum(["asc", "desc"]).default("asc"),
});

// API response wrapper schema
export const apiResponseSchema = <T extends z.ZodTypeAny>(dataSchema: T) =>
  z.object({
    success: z.boolean(),
    data: dataSchema.optional(),
    error: z.string().optional(),
    timestamp: z.date(),
  });

// Types
export type PaginationState = z.infer<typeof paginationStateSchema>;
export type PaginationResult = z.infer<typeof paginationResultSchema>;
export type DeviceFilters = z.infer<typeof deviceFiltersSchema>;
export type SecurityFilters = z.infer<typeof securityFiltersSchema>;
export type SortOptions = z.infer<typeof sortOptionsSchema>;
export type ApiResponse<T> = {
  success: boolean;
  data?: T;
  error?: string;
  timestamp: Date;
};

// Paginated response type
export interface PaginatedResponse<T> {
  items: T[];
  pagination: PaginationResult;
}

// Sort direction enum
export const SortDirection = {
  ASC: "asc",
  DESC: "desc",
} as const;

export type SortDirectionType =
  (typeof SortDirection)[keyof typeof SortDirection];

// Common field validation patterns
export const ValidationPatterns = {
  IP_ADDRESS:
    /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/,
  MAC_ADDRESS: /^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$/,
  HOSTNAME:
    /^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/,
  CRON_EXPRESSION:
    /^(\*|([0-9]|1[0-9]|2[0-9]|3[0-9]|4[0-9]|5[0-9])|\*\/([0-9]|1[0-9]|2[0-9]|3[0-9]|4[0-9]|5[0-9])) (\*|([0-9]|1[0-9]|2[0-3])|\*\/([0-9]|1[0-9]|2[0-3])) (\*|([1-9]|1[0-9]|2[0-9]|3[0-1])|\*\/([1-9]|1[0-9]|2[0-9]|3[0-1])) (\*|([1-9]|1[0-2])|\*\/([1-9]|1[0-2])) (\*|([0-6])|\*\/([0-6]))$/,
} as const;

// Error types
export interface ValidationError {
  field: string;
  message: string;
  code: string;
}

export interface AppError {
  code: string;
  message: string;
  details?: string;
  timestamp: Date;
}

// Loading states
export interface LoadingState {
  isLoading: boolean;
  error?: string;
  lastUpdated?: Date;
}
