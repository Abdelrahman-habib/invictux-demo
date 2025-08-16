import { z } from "zod";

// Security check result validation schema
export const securityCheckResultSchema = z.object({
  id: z.string(),
  deviceId: z.string(),
  checkName: z.string().min(1, "Check name is required"),
  checkType: z.string().min(1, "Check type is required"),
  severity: z.enum(["critical", "high", "medium", "low", "info"], {
    message: "Please select a valid severity level",
  }),
  status: z.enum(["pass", "fail", "warning", "error"], {
    message: "Please select a valid status",
  }),
  message: z.string().optional(),
  evidence: z.string().optional(),
  checkedAt: z.date(),
});

// Security rule validation schema
export const securityRuleSchema = z.object({
  id: z.string(),
  name: z.string().min(1, "Rule name is required"),
  description: z.string().optional(),
  vendor: z.enum(["cisco", "juniper", "hp", "aruba"], {
    message: "Please select a valid vendor",
  }),
  command: z.string().min(1, "Command is required"),
  expectedPattern: z.string().optional(),
  severity: z.enum(["critical", "high", "medium", "low", "info"], {
    message: "Please select a valid severity level",
  }),
  enabled: z.boolean().default(true),
  createdAt: z.date(),
});

// Security check execution request schema
export const securityCheckRequestSchema = z.object({
  deviceIds: z.array(z.string()).min(1, "At least one device must be selected"),
  ruleIds: z.array(z.string()).optional(),
  parallel: z.boolean().default(true),
});

// Types
export type SecurityCheckResult = z.infer<typeof securityCheckResultSchema>;
export type SecurityRule = z.infer<typeof securityRuleSchema>;
export type SecurityCheckRequest = z.infer<typeof securityCheckRequestSchema>;

// Security check status enum
export const SecurityCheckStatus = {
  PASS: "pass",
  FAIL: "fail",
  WARNING: "warning",
  ERROR: "error",
} as const;

export type SecurityCheckStatusType =
  (typeof SecurityCheckStatus)[keyof typeof SecurityCheckStatus];

// Severity enum
export const Severity = {
  CRITICAL: "critical",
  HIGH: "high",
  MEDIUM: "medium",
  LOW: "low",
  INFO: "info",
} as const;

export type SeverityType = (typeof Severity)[keyof typeof Severity];

// Security check summary interface
export interface SecurityCheckSummary {
  deviceId: string;
  deviceName: string;
  totalChecks: number;
  passedChecks: number;
  failedChecks: number;
  warningChecks: number;
  errorChecks: number;
  lastChecked?: Date;
  overallStatus: SecurityCheckStatusType;
}
