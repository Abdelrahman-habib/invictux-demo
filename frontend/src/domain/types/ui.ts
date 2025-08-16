import { z } from "zod";

// UI state schemas for different features
export const deviceUIStateSchema = z.object({
  currentPage: z.number().min(1).default(1),
  searchQuery: z.string().default(""),
  filters: z
    .object({
      deviceType: z.string().optional(),
      vendor: z.string().optional(),
      status: z.string().optional(),
      tags: z.array(z.string()).optional(),
    })
    .default({}),
  sortBy: z.string().default("name"),
  sortDirection: z.enum(["asc", "desc"]).default("asc"),
  isSelectionMode: z.boolean().default(false),
  selectedDeviceIds: z.set(z.string()).default(new Set()),
  isFormModalOpen: z.boolean().default(false),
  editingDevice: z.string().optional(), // Device ID when editing
});

export const securityUIStateSchema = z.object({
  currentPage: z.number().min(1).default(1),
  searchQuery: z.string().default(""),
  filters: z
    .object({
      severity: z.string().optional(),
      status: z.string().optional(),
      deviceId: z.string().optional(),
      checkType: z.string().optional(),
    })
    .default({}),
  sortBy: z.string().default("checkedAt"),
  sortDirection: z.enum(["asc", "desc"]).default("desc"),
  selectedIssueIds: z.set(z.string()).default(new Set()),
  isRunningChecks: z.boolean().default(false),
  checkProgress: z.number().min(0).max(100).default(0),
});

export const reportUIStateSchema = z.object({
  currentPage: z.number().min(1).default(1),
  searchQuery: z.string().default(""),
  filters: z
    .object({
      type: z.string().optional(),
      status: z.string().optional(),
      format: z.string().optional(),
    })
    .default({}),
  sortBy: z.string().default("createdAt"),
  sortDirection: z.enum(["asc", "desc"]).default("desc"),
  isGenerating: z.boolean().default(false),
  generationProgress: z.number().min(0).max(100).default(0),
  isScheduleModalOpen: z.boolean().default(false),
});

// Types
export type DeviceUIState = z.infer<typeof deviceUIStateSchema>;
export type SecurityUIState = z.infer<typeof securityUIStateSchema>;
export type ReportUIState = z.infer<typeof reportUIStateSchema>;

// Modal states
export interface ModalState<T = unknown> {
  isOpen: boolean;
  data?: T;
}

export interface FormModalState<T = unknown> extends ModalState<T> {
  mode: "create" | "edit";
  initialData?: T;
}

// Selection state
export interface SelectionState<T = string> {
  selectedIds: Set<T>;
  isSelectionMode: boolean;
  selectAll: boolean;
}

// Progress state
export interface ProgressState {
  isActive: boolean;
  progress: number;
  message?: string;
  estimatedTimeRemaining?: number;
}

// Toast/notification types
export interface ToastMessage {
  id: string;
  type: "success" | "error" | "warning" | "info";
  title: string;
  message?: string;
  duration?: number;
  action?: {
    label: string;
    onClick: () => void;
  };
}

// Theme and appearance
export interface ThemeState {
  mode: "light" | "dark" | "system";
  primaryColor: string;
  fontSize: "small" | "medium" | "large";
  compactMode: boolean;
}

// Navigation state
export interface NavigationState {
  currentRoute: string;
  breadcrumbs: Array<{
    label: string;
    path: string;
  }>;
  sidebarCollapsed: boolean;
}

// Search and filter state
export interface SearchState {
  query: string;
  suggestions: string[];
  recentSearches: string[];
  isSearching: boolean;
}

export interface FilterState<T = Record<string, unknown>> {
  activeFilters: T;
  availableFilters: Array<{
    key: keyof T;
    label: string;
    type: "select" | "multiselect" | "date" | "range";
    options?: Array<{ value: string; label: string }>;
  }>;
  filterCount: number;
}

// Table/list view state
export interface TableState {
  columns: Array<{
    key: string;
    label: string;
    sortable: boolean;
    visible: boolean;
    width?: number;
  }>;
  sortBy: string;
  sortDirection: "asc" | "desc";
  pageSize: number;
  currentPage: number;
  selectedRows: Set<string>;
  expandedRows: Set<string>;
}
