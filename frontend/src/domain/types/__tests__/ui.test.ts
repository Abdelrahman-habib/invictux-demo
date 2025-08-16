import { describe, it, expect } from "vitest";
import {
  deviceUIStateSchema,
  securityUIStateSchema,
  reportUIStateSchema,
  type ModalState,
  type FormModalState,
  type SelectionState,
  type ProgressState,
  type ToastMessage,
  type ThemeState,
  type NavigationState,
  type SearchState,
  type FilterState,
  type TableState,
} from "../ui";

describe("UI Types", () => {
  describe("deviceUIStateSchema validation", () => {
    const validDeviceUIState = {
      currentPage: 1,
      searchQuery: "test",
      filters: {
        deviceType: "router",
        vendor: "cisco",
        status: "online",
        tags: ["production"],
      },
      sortBy: "name",
      sortDirection: "asc" as const,
      isSelectionMode: false,
      selectedDeviceIds: new Set(["device-1", "device-2"]),
      isFormModalOpen: false,
      editingDevice: "device-123",
    };

    it("should validate valid device UI state", () => {
      const result = deviceUIStateSchema.safeParse(validDeviceUIState);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data.currentPage).toBe(1);
        expect(result.data.searchQuery).toBe("test");
        expect(result.data.sortBy).toBe("name");
        expect(result.data.sortDirection).toBe("asc");
      }
    });

    it("should use default values when not provided", () => {
      const result = deviceUIStateSchema.safeParse({});
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data.currentPage).toBe(1);
        expect(result.data.searchQuery).toBe("");
        expect(result.data.sortBy).toBe("name");
        expect(result.data.sortDirection).toBe("asc");
        expect(result.data.isSelectionMode).toBe(false);
        expect(result.data.isFormModalOpen).toBe(false);
      }
    });

    it("should accept both asc and desc sort directions", () => {
      const directions = ["asc", "desc"] as const;
      directions.forEach((direction) => {
        const state = { ...validDeviceUIState, sortDirection: direction };
        const result = deviceUIStateSchema.safeParse(state);
        expect(result.success).toBe(true);
      });
    });

    it("should reject invalid page numbers", () => {
      const invalidState = { ...validDeviceUIState, currentPage: 0 };
      const result = deviceUIStateSchema.safeParse(invalidState);
      expect(result.success).toBe(false);
    });
  });

  describe("securityUIStateSchema validation", () => {
    const validSecurityUIState = {
      currentPage: 1,
      searchQuery: "security",
      filters: {
        severity: "high",
        status: "fail",
        deviceId: "device-123",
        checkType: "authentication",
      },
      sortBy: "checkedAt",
      sortDirection: "desc" as const,
      selectedIssueIds: new Set(["issue-1", "issue-2"]),
      isRunningChecks: false,
      checkProgress: 75,
    };

    it("should validate valid security UI state", () => {
      const result = securityUIStateSchema.safeParse(validSecurityUIState);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data.sortBy).toBe("checkedAt");
        expect(result.data.sortDirection).toBe("desc");
        expect(result.data.checkProgress).toBe(75);
      }
    });

    it("should use default values when not provided", () => {
      const result = securityUIStateSchema.safeParse({});
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data.currentPage).toBe(1);
        expect(result.data.searchQuery).toBe("");
        expect(result.data.sortBy).toBe("checkedAt");
        expect(result.data.sortDirection).toBe("desc");
        expect(result.data.isRunningChecks).toBe(false);
        expect(result.data.checkProgress).toBe(0);
      }
    });

    it("should reject invalid progress values", () => {
      const invalidStates = [
        { ...validSecurityUIState, checkProgress: -1 },
        { ...validSecurityUIState, checkProgress: 101 },
      ];

      invalidStates.forEach((state) => {
        const result = securityUIStateSchema.safeParse(state);
        expect(result.success).toBe(false);
      });
    });

    it("should accept valid progress values", () => {
      const validProgress = [0, 25, 50, 75, 100];
      validProgress.forEach((progress) => {
        const state = { ...validSecurityUIState, checkProgress: progress };
        const result = securityUIStateSchema.safeParse(state);
        expect(result.success).toBe(true);
      });
    });
  });

  describe("reportUIStateSchema validation", () => {
    const validReportUIState = {
      currentPage: 1,
      searchQuery: "report",
      filters: {
        type: "executive",
        status: "completed",
        format: "pdf",
      },
      sortBy: "createdAt",
      sortDirection: "desc" as const,
      isGenerating: false,
      generationProgress: 50,
      isScheduleModalOpen: false,
    };

    it("should validate valid report UI state", () => {
      const result = reportUIStateSchema.safeParse(validReportUIState);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data.sortBy).toBe("createdAt");
        expect(result.data.generationProgress).toBe(50);
        expect(result.data.isGenerating).toBe(false);
      }
    });

    it("should use default values when not provided", () => {
      const result = reportUIStateSchema.safeParse({});
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data.currentPage).toBe(1);
        expect(result.data.searchQuery).toBe("");
        expect(result.data.sortBy).toBe("createdAt");
        expect(result.data.sortDirection).toBe("desc");
        expect(result.data.isGenerating).toBe(false);
        expect(result.data.generationProgress).toBe(0);
        expect(result.data.isScheduleModalOpen).toBe(false);
      }
    });
  });

  describe("Interface types", () => {
    it("should have correct ModalState structure", () => {
      const modalState: ModalState<string> = {
        isOpen: true,
        data: "test data",
      };

      expect(modalState.isOpen).toBe(true);
      expect(modalState.data).toBe("test data");
    });

    it("should have correct FormModalState structure", () => {
      const formModalState: FormModalState<{ name: string }> = {
        isOpen: true,
        mode: "edit",
        data: { name: "test" },
        initialData: { name: "initial" },
      };

      expect(formModalState.isOpen).toBe(true);
      expect(formModalState.mode).toBe("edit");
      expect(formModalState.data?.name).toBe("test");
      expect(formModalState.initialData?.name).toBe("initial");
    });

    it("should have correct SelectionState structure", () => {
      const selectionState: SelectionState<string> = {
        selectedIds: new Set(["id1", "id2"]),
        isSelectionMode: true,
        selectAll: false,
      };

      expect(selectionState.selectedIds.has("id1")).toBe(true);
      expect(selectionState.isSelectionMode).toBe(true);
      expect(selectionState.selectAll).toBe(false);
    });

    it("should have correct ProgressState structure", () => {
      const progressState: ProgressState = {
        isActive: true,
        progress: 75,
        message: "Processing...",
        estimatedTimeRemaining: 30,
      };

      expect(progressState.isActive).toBe(true);
      expect(progressState.progress).toBe(75);
      expect(progressState.message).toBe("Processing...");
      expect(progressState.estimatedTimeRemaining).toBe(30);
    });

    it("should have correct ToastMessage structure", () => {
      const toastMessage: ToastMessage = {
        id: "toast-1",
        type: "success",
        title: "Success",
        message: "Operation completed",
        duration: 5000,
        action: {
          label: "Undo",
          onClick: () => console.log("Undo clicked"),
        },
      };

      expect(toastMessage.id).toBe("toast-1");
      expect(toastMessage.type).toBe("success");
      expect(toastMessage.title).toBe("Success");
      expect(toastMessage.message).toBe("Operation completed");
      expect(toastMessage.duration).toBe(5000);
      expect(toastMessage.action?.label).toBe("Undo");
    });

    it("should have correct ThemeState structure", () => {
      const themeState: ThemeState = {
        mode: "dark",
        primaryColor: "#3b82f6",
        fontSize: "medium",
        compactMode: false,
      };

      expect(themeState.mode).toBe("dark");
      expect(themeState.primaryColor).toBe("#3b82f6");
      expect(themeState.fontSize).toBe("medium");
      expect(themeState.compactMode).toBe(false);
    });

    it("should have correct NavigationState structure", () => {
      const navigationState: NavigationState = {
        currentRoute: "/devices",
        breadcrumbs: [
          { label: "Home", path: "/" },
          { label: "Devices", path: "/devices" },
        ],
        sidebarCollapsed: false,
      };

      expect(navigationState.currentRoute).toBe("/devices");
      expect(navigationState.breadcrumbs).toHaveLength(2);
      expect(navigationState.sidebarCollapsed).toBe(false);
    });

    it("should have correct SearchState structure", () => {
      const searchState: SearchState = {
        query: "test query",
        suggestions: ["suggestion1", "suggestion2"],
        recentSearches: ["recent1", "recent2"],
        isSearching: false,
      };

      expect(searchState.query).toBe("test query");
      expect(searchState.suggestions).toHaveLength(2);
      expect(searchState.recentSearches).toHaveLength(2);
      expect(searchState.isSearching).toBe(false);
    });

    it("should have correct FilterState structure", () => {
      const filterState: FilterState<{ status: string; type: string }> = {
        activeFilters: { status: "active", type: "router" },
        availableFilters: [
          {
            key: "status",
            label: "Status",
            type: "select",
            options: [
              { value: "active", label: "Active" },
              { value: "inactive", label: "Inactive" },
            ],
          },
          {
            key: "type",
            label: "Type",
            type: "multiselect",
          },
        ],
        filterCount: 2,
      };

      expect(filterState.activeFilters.status).toBe("active");
      expect(filterState.availableFilters).toHaveLength(2);
      expect(filterState.filterCount).toBe(2);
    });

    it("should have correct TableState structure", () => {
      const tableState: TableState = {
        columns: [
          {
            key: "name",
            label: "Name",
            sortable: true,
            visible: true,
            width: 200,
          },
          {
            key: "status",
            label: "Status",
            sortable: false,
            visible: true,
          },
        ],
        sortBy: "name",
        sortDirection: "asc",
        pageSize: 20,
        currentPage: 1,
        selectedRows: new Set(["row1", "row2"]),
        expandedRows: new Set(["row1"]),
      };

      expect(tableState.columns).toHaveLength(2);
      expect(tableState.sortBy).toBe("name");
      expect(tableState.sortDirection).toBe("asc");
      expect(tableState.pageSize).toBe(20);
      expect(tableState.selectedRows.has("row1")).toBe(true);
      expect(tableState.expandedRows.has("row1")).toBe(true);
    });
  });
});
