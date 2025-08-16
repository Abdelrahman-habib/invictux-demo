import { describe, it, expect } from "vitest";
import {
  paginationStateSchema,
  paginationResultSchema,
  deviceFiltersSchema,
  securityFiltersSchema,
  sortOptionsSchema,
  apiResponseSchema,
  ValidationPatterns,
  SortDirection,
  type PaginatedResponse,
  type ValidationError,
  type AppError,
  type LoadingState,
} from "../common";

describe("Common Types", () => {
  describe("paginationStateSchema validation", () => {
    const validPagination = {
      currentPage: 1,
      limit: 20,
      searchQuery: "test query",
    };

    it("should validate valid pagination state", () => {
      const result = paginationStateSchema.safeParse(validPagination);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data).toEqual(validPagination);
      }
    });

    it("should use default values when not provided", () => {
      const result = paginationStateSchema.safeParse({});
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data.currentPage).toBe(1);
        expect(result.data.limit).toBe(20);
        expect(result.data.searchQuery).toBe("");
      }
    });

    it("should reject invalid page numbers", () => {
      const invalid = { ...validPagination, currentPage: 0 };
      const result = paginationStateSchema.safeParse(invalid);
      expect(result.success).toBe(false);
    });

    it("should reject invalid limits", () => {
      const invalidLimits = [0, -1, 101];
      invalidLimits.forEach((limit) => {
        const invalid = { ...validPagination, limit };
        const result = paginationStateSchema.safeParse(invalid);
        expect(result.success).toBe(false);
      });
    });

    it("should accept valid limits", () => {
      const validLimits = [1, 20, 50, 100];
      validLimits.forEach((limit) => {
        const valid = { ...validPagination, limit };
        const result = paginationStateSchema.safeParse(valid);
        expect(result.success).toBe(true);
      });
    });
  });

  describe("paginationResultSchema validation", () => {
    const validResult = {
      totalCount: 100,
      currentPage: 1,
      totalPages: 5,
      hasNextPage: true,
      hasPrevPage: false,
    };

    it("should validate valid pagination result", () => {
      const result = paginationResultSchema.safeParse(validResult);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data).toEqual(validResult);
      }
    });
  });

  describe("deviceFiltersSchema validation", () => {
    const validFilters = {
      deviceType: "router",
      vendor: "cisco",
      status: "online",
      tags: ["production", "core"],
    };

    it("should validate valid device filters", () => {
      const result = deviceFiltersSchema.safeParse(validFilters);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data).toEqual(validFilters);
      }
    });

    it("should allow all fields to be optional", () => {
      const result = deviceFiltersSchema.safeParse({});
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data.deviceType).toBeUndefined();
        expect(result.data.vendor).toBeUndefined();
        expect(result.data.status).toBeUndefined();
        expect(result.data.tags).toBeUndefined();
      }
    });

    it("should allow partial filters", () => {
      const partialFilters = {
        deviceType: "switch",
        status: "offline",
      };
      const result = deviceFiltersSchema.safeParse(partialFilters);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data.deviceType).toBe("switch");
        expect(result.data.status).toBe("offline");
        expect(result.data.vendor).toBeUndefined();
        expect(result.data.tags).toBeUndefined();
      }
    });
  });

  describe("securityFiltersSchema validation", () => {
    const validFilters = {
      severity: "high",
      status: "fail",
      deviceId: "device-123",
      checkType: "authentication",
      dateRange: {
        startDate: new Date("2024-01-01"),
        endDate: new Date("2024-01-31"),
      },
    };

    it("should validate valid security filters", () => {
      const result = securityFiltersSchema.safeParse(validFilters);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data).toEqual(validFilters);
      }
    });

    it("should allow all fields to be optional", () => {
      const result = securityFiltersSchema.safeParse({});
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data.severity).toBeUndefined();
        expect(result.data.status).toBeUndefined();
        expect(result.data.deviceId).toBeUndefined();
        expect(result.data.checkType).toBeUndefined();
        expect(result.data.dateRange).toBeUndefined();
      }
    });
  });

  describe("sortOptionsSchema validation", () => {
    const validSort = {
      field: "name",
      direction: "asc" as const,
    };

    it("should validate valid sort options", () => {
      const result = sortOptionsSchema.safeParse(validSort);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data).toEqual(validSort);
      }
    });

    it("should use default direction when not provided", () => {
      const sortWithoutDirection = { field: "createdAt" };
      const result = sortOptionsSchema.safeParse(sortWithoutDirection);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data.direction).toBe("asc");
      }
    });

    it("should accept both asc and desc directions", () => {
      const directions = ["asc", "desc"] as const;
      directions.forEach((direction) => {
        const valid = { ...validSort, direction };
        const result = sortOptionsSchema.safeParse(valid);
        expect(result.success).toBe(true);
      });
    });
  });

  describe("apiResponseSchema validation", () => {
    it("should validate successful API response", () => {
      const stringSchema = apiResponseSchema(
        paginationStateSchema.pick({ currentPage: true })
      );
      const validResponse = {
        success: true,
        data: { currentPage: 1 },
        timestamp: new Date(),
      };

      const result = stringSchema.safeParse(validResponse);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data.success).toBe(true);
        expect(result.data.data).toEqual({ currentPage: 1 });
      }
    });

    it("should validate error API response", () => {
      const stringSchema = apiResponseSchema(
        paginationStateSchema.pick({ currentPage: true })
      );
      const errorResponse = {
        success: false,
        error: "Something went wrong",
        timestamp: new Date(),
      };

      const result = stringSchema.safeParse(errorResponse);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data.success).toBe(false);
        expect(result.data.error).toBe("Something went wrong");
        expect(result.data.data).toBeUndefined();
      }
    });
  });

  describe("ValidationPatterns", () => {
    it("should validate IP addresses correctly", () => {
      const validIPs = [
        "0.0.0.0",
        "192.168.1.1",
        "10.0.0.1",
        "172.16.0.1",
        "255.255.255.255",
      ];
      const invalidIPs = [
        "256.1.1.1",
        "192.168.1",
        "192.168.1.1.1",
        "not-an-ip",
        "192.168.1.256",
      ];

      validIPs.forEach((ip) => {
        expect(ValidationPatterns.IP_ADDRESS.test(ip)).toBe(true);
      });

      invalidIPs.forEach((ip) => {
        expect(ValidationPatterns.IP_ADDRESS.test(ip)).toBe(false);
      });
    });

    it("should validate MAC addresses correctly", () => {
      const validMACs = [
        "00:11:22:33:44:55",
        "AA:BB:CC:DD:EE:FF",
        "00-11-22-33-44-55",
        "aa:bb:cc:dd:ee:ff",
      ];
      const invalidMACs = [
        "00:11:22:33:44",
        "00:11:22:33:44:55:66",
        "GG:HH:II:JJ:KK:LL",
        "00.11.22.33.44.55",
      ];

      validMACs.forEach((mac) => {
        expect(ValidationPatterns.MAC_ADDRESS.test(mac)).toBe(true);
      });

      invalidMACs.forEach((mac) => {
        expect(ValidationPatterns.MAC_ADDRESS.test(mac)).toBe(false);
      });
    });

    it("should validate hostnames correctly", () => {
      const validHostnames = [
        "example.com",
        "sub.example.com",
        "test-server",
        "server1",
        "a.b.c.d",
      ];
      const invalidHostnames = [
        "-invalid",
        "invalid-",
        "invalid..double.dot",
        "",
        ".invalid",
        "invalid.",
      ];

      validHostnames.forEach((hostname) => {
        expect(ValidationPatterns.HOSTNAME.test(hostname)).toBe(true);
      });

      invalidHostnames.forEach((hostname) => {
        expect(ValidationPatterns.HOSTNAME.test(hostname)).toBe(false);
      });
    });

    it("should validate cron expressions correctly", () => {
      const validCrons = [
        "0 0 * * *", // Daily at midnight
        "0 9 * * 1", // Weekly on Monday at 9 AM
        "*/15 * * * *", // Every 15 minutes
        "0 0 1 * *", // Monthly on the 1st
        "0 0 * * 0", // Weekly on Sunday
      ];
      const invalidCrons = [
        "60 0 * * *", // Invalid minute
        "0 25 * * *", // Invalid hour
        "0 0 32 * *", // Invalid day
        "0 0 * 13 *", // Invalid month
        "0 0 * * 7", // Invalid day of week
        "invalid cron",
      ];

      validCrons.forEach((cron) => {
        expect(ValidationPatterns.CRON_EXPRESSION.test(cron)).toBe(true);
      });

      invalidCrons.forEach((cron) => {
        expect(ValidationPatterns.CRON_EXPRESSION.test(cron)).toBe(false);
      });
    });
  });

  describe("SortDirection enum", () => {
    it("should have correct values", () => {
      expect(SortDirection.ASC).toBe("asc");
      expect(SortDirection.DESC).toBe("desc");
    });
  });

  describe("Type interfaces", () => {
    it("should have correct PaginatedResponse structure", () => {
      const response: PaginatedResponse<string> = {
        items: ["item1", "item2"],
        pagination: {
          totalCount: 2,
          currentPage: 1,
          totalPages: 1,
          hasNextPage: false,
          hasPrevPage: false,
        },
      };

      expect(response.items).toHaveLength(2);
      expect(response.pagination.totalCount).toBe(2);
    });

    it("should have correct ValidationError structure", () => {
      const error: ValidationError = {
        field: "email",
        message: "Invalid email format",
        code: "INVALID_EMAIL",
      };

      expect(error.field).toBe("email");
      expect(error.message).toBe("Invalid email format");
      expect(error.code).toBe("INVALID_EMAIL");
    });

    it("should have correct AppError structure", () => {
      const error: AppError = {
        code: "NETWORK_ERROR",
        message: "Failed to connect to device",
        details: "Connection timeout after 30 seconds",
        timestamp: new Date(),
      };

      expect(error.code).toBe("NETWORK_ERROR");
      expect(error.message).toBe("Failed to connect to device");
      expect(error.details).toBe("Connection timeout after 30 seconds");
      expect(error.timestamp).toBeInstanceOf(Date);
    });

    it("should have correct LoadingState structure", () => {
      const loadingState: LoadingState = {
        isLoading: true,
        error: "Failed to load data",
        lastUpdated: new Date(),
      };

      expect(loadingState.isLoading).toBe(true);
      expect(loadingState.error).toBe("Failed to load data");
      expect(loadingState.lastUpdated).toBeInstanceOf(Date);
    });
  });
});
