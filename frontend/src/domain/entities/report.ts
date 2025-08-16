import { z } from "zod";

// Report generation request schema
export const reportRequestSchema = z.object({
  name: z.string().min(1, "Report name is required"),
  type: z.enum(["executive", "technical", "compliance"], {
    message: "Please select a valid report type",
  }),
  format: z.enum(["pdf", "csv"], {
    message: "Please select a valid format",
  }),
  deviceIds: z.array(z.string()).optional(),
  dateRange: z
    .object({
      startDate: z.date(),
      endDate: z.date(),
    })
    .refine((data) => data.startDate <= data.endDate, {
      message: "Start date must be before or equal to end date",
      path: ["endDate"],
    }),
  includePassedChecks: z.boolean().default(false),
  includeSummary: z.boolean().default(true),
  includeRecommendations: z.boolean().default(true),
});

// Report metadata schema
export const reportMetadataSchema = z.object({
  id: z.string(),
  name: z.string(),
  type: z.enum(["executive", "technical", "compliance"]),
  format: z.enum(["pdf", "csv"]),
  status: z.enum(["pending", "generating", "completed", "failed"]),
  filePath: z.string().optional(),
  fileSize: z.number().optional(),
  deviceCount: z.number(),
  checkCount: z.number(),
  generatedAt: z.date().optional(),
  createdAt: z.date(),
  error: z.string().optional(),
});

// Report schedule schema
export const reportScheduleSchema = z.object({
  id: z.string(),
  name: z.string().min(1, "Schedule name is required"),
  reportRequest: reportRequestSchema.omit({ name: true }),
  cronExpression: z.string().min(1, "Cron expression is required"),
  emailRecipients: z
    .array(z.string().email("Invalid email address"))
    .optional(),
  enabled: z.boolean().default(true),
  lastRun: z.date().optional(),
  nextRun: z.date().optional(),
  createdAt: z.date(),
  updatedAt: z.date(),
});

// Types
export type ReportRequest = z.infer<typeof reportRequestSchema>;
export type ReportMetadata = z.infer<typeof reportMetadataSchema>;
export type ReportSchedule = z.infer<typeof reportScheduleSchema>;

// Report type enum
export const ReportType = {
  EXECUTIVE: "executive",
  TECHNICAL: "technical",
  COMPLIANCE: "compliance",
} as const;

export type ReportTypeType = (typeof ReportType)[keyof typeof ReportType];

// Report format enum
export const ReportFormat = {
  PDF: "pdf",
  CSV: "csv",
} as const;

export type ReportFormatType = (typeof ReportFormat)[keyof typeof ReportFormat];

// Report status enum
export const ReportStatus = {
  PENDING: "pending",
  GENERATING: "generating",
  COMPLETED: "completed",
  FAILED: "failed",
} as const;

export type ReportStatusType = (typeof ReportStatus)[keyof typeof ReportStatus];

// Report summary interface
export interface ReportSummary {
  totalReports: number;
  pendingReports: number;
  completedReports: number;
  failedReports: number;
  totalFileSize: number;
  lastGenerated?: Date;
}
