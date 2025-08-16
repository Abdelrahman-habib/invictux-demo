import { describe, it, expect } from "vitest";
import {
  reportRequestSchema,
  reportMetadataSchema,
  reportScheduleSchema,
  ReportType,
  ReportFormat,
  ReportStatus,
  type ReportSummary,
} from "../report";

describe("Report Entity", () => {
  describe("reportRequestSchema validation", () => {
    const validRequest = {
      name: "Monthly Security Report",
      type: "executive" as const,
      format: "pdf" as const,
      deviceIds: ["device-1", "device-2"],
      dateRange: {
        startDate: new Date("2024-01-01"),
        endDate: new Date("2024-01-31"),
      },
      includePassedChecks: false,
      includeSummary: true,
      includeRecommendations: true,
    };

    it("should validate a valid report request", () => {
      const result = reportRequestSchema.safeParse(validRequest);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data).toEqual(validRequest);
      }
    });

    it("should use default values for optional boolean fields", () => {
      const minimalRequest = {
        name: "Test Report",
        type: "technical" as const,
        format: "csv" as const,
        dateRange: {
          startDate: new Date("2024-01-01"),
          endDate: new Date("2024-01-31"),
        },
      };

      const result = reportRequestSchema.safeParse(minimalRequest);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data.includePassedChecks).toBe(false);
        expect(result.data.includeSummary).toBe(true);
        expect(result.data.includeRecommendations).toBe(true);
      }
    });

    describe("field validation", () => {
      it("should reject empty report name", () => {
        const invalid = { ...validRequest, name: "" };
        const result = reportRequestSchema.safeParse(invalid);

        expect(result.success).toBe(false);
        if (!result.success) {
          expect(result.error.issues[0].message).toBe(
            "Report name is required"
          );
        }
      });

      it("should reject invalid report type", () => {
        const invalid = { ...validRequest, type: "invalid" as never };
        const result = reportRequestSchema.safeParse(invalid);

        expect(result.success).toBe(false);
        if (!result.success) {
          expect(result.error.issues[0].message).toBe(
            "Please select a valid report type"
          );
        }
      });

      it("should accept all valid report types", () => {
        const validTypes = ["executive", "technical", "compliance"] as const;

        validTypes.forEach((type) => {
          const valid = { ...validRequest, type };
          const result = reportRequestSchema.safeParse(valid);
          expect(result.success).toBe(true);
        });
      });

      it("should reject invalid format", () => {
        const invalid = { ...validRequest, format: "invalid" as never };
        const result = reportRequestSchema.safeParse(invalid);

        expect(result.success).toBe(false);
        if (!result.success) {
          expect(result.error.issues[0].message).toBe(
            "Please select a valid format"
          );
        }
      });

      it("should accept all valid formats", () => {
        const validFormats = ["pdf", "csv"] as const;

        validFormats.forEach((format) => {
          const valid = { ...validRequest, format };
          const result = reportRequestSchema.safeParse(valid);
          expect(result.success).toBe(true);
        });
      });

      it("should reject invalid date range (start after end)", () => {
        const invalid = {
          ...validRequest,
          dateRange: {
            startDate: new Date("2024-01-31"),
            endDate: new Date("2024-01-01"),
          },
        };
        const result = reportRequestSchema.safeParse(invalid);

        expect(result.success).toBe(false);
        if (!result.success) {
          expect(result.error.issues[0].message).toBe(
            "Start date must be before or equal to end date"
          );
          expect(result.error.issues[0].path).toEqual(["dateRange", "endDate"]);
        }
      });

      it("should accept equal start and end dates", () => {
        const sameDate = new Date("2024-01-15");
        const valid = {
          ...validRequest,
          dateRange: {
            startDate: sameDate,
            endDate: sameDate,
          },
        };
        const result = reportRequestSchema.safeParse(valid);
        expect(result.success).toBe(true);
      });
    });
  });

  describe("reportMetadataSchema validation", () => {
    const validMetadata = {
      id: "report-123",
      name: "Security Report",
      type: "executive" as const,
      format: "pdf" as const,
      status: "completed" as const,
      filePath: "/reports/security-report.pdf",
      fileSize: 1024000,
      deviceCount: 5,
      checkCount: 50,
      generatedAt: new Date(),
      createdAt: new Date(),
    };

    it("should validate valid report metadata", () => {
      const result = reportMetadataSchema.safeParse(validMetadata);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data).toEqual(validMetadata);
      }
    });

    it("should allow optional fields to be undefined", () => {
      const minimalMetadata = {
        id: "report-456",
        name: "Test Report",
        type: "technical" as const,
        format: "csv" as const,
        status: "pending" as const,
        deviceCount: 3,
        checkCount: 30,
        createdAt: new Date(),
      };

      const result = reportMetadataSchema.safeParse(minimalMetadata);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data.filePath).toBeUndefined();
        expect(result.data.fileSize).toBeUndefined();
        expect(result.data.generatedAt).toBeUndefined();
        expect(result.data.error).toBeUndefined();
      }
    });

    it("should accept all valid status values", () => {
      const validStatuses = [
        "pending",
        "generating",
        "completed",
        "failed",
      ] as const;

      validStatuses.forEach((status) => {
        const valid = { ...validMetadata, status };
        const result = reportMetadataSchema.safeParse(valid);
        expect(result.success).toBe(true);
      });
    });
  });

  describe("reportScheduleSchema validation", () => {
    const validSchedule = {
      id: "schedule-123",
      name: "Weekly Security Report",
      reportRequest: {
        type: "executive" as const,
        format: "pdf" as const,
        dateRange: {
          startDate: new Date("2024-01-01"),
          endDate: new Date("2024-01-07"),
        },
        includePassedChecks: false,
        includeSummary: true,
        includeRecommendations: true,
      },
      cronExpression: "0 9 * * 1",
      emailRecipients: ["admin@example.com", "security@example.com"],
      enabled: true,
      lastRun: new Date(),
      nextRun: new Date(),
      createdAt: new Date(),
      updatedAt: new Date(),
    };

    it("should validate a valid report schedule", () => {
      const result = reportScheduleSchema.safeParse(validSchedule);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data).toEqual(validSchedule);
      }
    });

    it("should use default enabled value", () => {
      const scheduleData = {
        id: validSchedule.id,
        name: validSchedule.name,
        reportRequest: validSchedule.reportRequest,
        cronExpression: validSchedule.cronExpression,
        emailRecipients: validSchedule.emailRecipients,
        lastRun: validSchedule.lastRun,
        nextRun: validSchedule.nextRun,
        createdAt: validSchedule.createdAt,
        updatedAt: validSchedule.updatedAt,
      };

      const result = reportScheduleSchema.safeParse(scheduleData);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data.enabled).toBe(true);
      }
    });

    it("should allow optional fields to be undefined", () => {
      const minimalSchedule = {
        id: "schedule-456",
        name: "Simple Schedule",
        reportRequest: {
          type: "technical" as const,
          format: "csv" as const,
          dateRange: {
            startDate: new Date("2024-01-01"),
            endDate: new Date("2024-01-31"),
          },
          includePassedChecks: false,
          includeSummary: true,
          includeRecommendations: true,
        },
        cronExpression: "0 0 1 * *",
        createdAt: new Date(),
        updatedAt: new Date(),
      };

      const result = reportScheduleSchema.safeParse(minimalSchedule);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data.emailRecipients).toBeUndefined();
        expect(result.data.lastRun).toBeUndefined();
        expect(result.data.nextRun).toBeUndefined();
        expect(result.data.enabled).toBe(true);
      }
    });

    describe("field validation", () => {
      it("should reject empty schedule name", () => {
        const invalid = { ...validSchedule, name: "" };
        const result = reportScheduleSchema.safeParse(invalid);

        expect(result.success).toBe(false);
        if (!result.success) {
          expect(result.error.issues[0].message).toBe(
            "Schedule name is required"
          );
        }
      });

      it("should reject empty cron expression", () => {
        const invalid = { ...validSchedule, cronExpression: "" };
        const result = reportScheduleSchema.safeParse(invalid);

        expect(result.success).toBe(false);
        if (!result.success) {
          expect(result.error.issues[0].message).toBe(
            "Cron expression is required"
          );
        }
      });

      it("should reject invalid email addresses", () => {
        const invalid = {
          ...validSchedule,
          emailRecipients: ["invalid-email", "admin@example.com"],
        };
        const result = reportScheduleSchema.safeParse(invalid);

        expect(result.success).toBe(false);
        if (!result.success) {
          expect(result.error.issues[0].message).toBe("Invalid email address");
        }
      });

      it("should accept valid email addresses", () => {
        const validEmails = [
          "user@example.com",
          "admin@company.org",
          "test.user+tag@domain.co.uk",
        ];
        const valid = { ...validSchedule, emailRecipients: validEmails };
        const result = reportScheduleSchema.safeParse(valid);
        expect(result.success).toBe(true);
      });
    });
  });

  describe("Report enums", () => {
    it("should have correct ReportType values", () => {
      expect(ReportType.EXECUTIVE).toBe("executive");
      expect(ReportType.TECHNICAL).toBe("technical");
      expect(ReportType.COMPLIANCE).toBe("compliance");
    });

    it("should have correct ReportFormat values", () => {
      expect(ReportFormat.PDF).toBe("pdf");
      expect(ReportFormat.CSV).toBe("csv");
    });

    it("should have correct ReportStatus values", () => {
      expect(ReportStatus.PENDING).toBe("pending");
      expect(ReportStatus.GENERATING).toBe("generating");
      expect(ReportStatus.COMPLETED).toBe("completed");
      expect(ReportStatus.FAILED).toBe("failed");
    });
  });

  describe("ReportSummary interface", () => {
    it("should have correct structure", () => {
      const summary: ReportSummary = {
        totalReports: 25,
        pendingReports: 2,
        completedReports: 20,
        failedReports: 3,
        totalFileSize: 50000000,
        lastGenerated: new Date(),
      };

      expect(summary.totalReports).toBe(25);
      expect(summary.pendingReports).toBe(2);
      expect(summary.completedReports).toBe(20);
      expect(summary.failedReports).toBe(3);
      expect(summary.totalFileSize).toBe(50000000);
      expect(summary.lastGenerated).toBeDefined();
    });
  });
});
