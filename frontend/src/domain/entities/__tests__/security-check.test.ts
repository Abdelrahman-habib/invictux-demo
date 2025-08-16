import { describe, it, expect } from "vitest";
import {
  securityCheckResultSchema,
  securityRuleSchema,
  securityCheckRequestSchema,
  SecurityCheckStatus,
  Severity,
  type SecurityCheckSummary,
} from "../security-check";

describe("Security Check Entity", () => {
  describe("securityCheckResultSchema validation", () => {
    const validCheckResult = {
      id: "check-123",
      deviceId: "device-456",
      checkName: "Default Password Check",
      checkType: "authentication",
      severity: "high" as const,
      status: "fail" as const,
      message: "Default password detected",
      evidence: "show running-config | include username",
      checkedAt: new Date(),
    };

    it("should validate a valid security check result", () => {
      const result = securityCheckResultSchema.safeParse(validCheckResult);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data).toEqual(validCheckResult);
      }
    });

    it("should allow optional message and evidence", () => {
      const minimalResult = {
        id: "check-123",
        deviceId: "device-456",
        checkName: "Port Check",
        checkType: "configuration",
        severity: "low" as const,
        status: "pass" as const,
        checkedAt: new Date(),
      };

      const result = securityCheckResultSchema.safeParse(minimalResult);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data.message).toBeUndefined();
        expect(result.data.evidence).toBeUndefined();
      }
    });

    describe("field validation", () => {
      it("should reject empty check name", () => {
        const invalid = { ...validCheckResult, checkName: "" };
        const result = securityCheckResultSchema.safeParse(invalid);

        expect(result.success).toBe(false);
        if (!result.success) {
          expect(result.error.issues[0].message).toBe("Check name is required");
        }
      });

      it("should reject empty check type", () => {
        const invalid = { ...validCheckResult, checkType: "" };
        const result = securityCheckResultSchema.safeParse(invalid);

        expect(result.success).toBe(false);
        if (!result.success) {
          expect(result.error.issues[0].message).toBe("Check type is required");
        }
      });

      it("should reject invalid severity levels", () => {
        const invalid = { ...validCheckResult, severity: "invalid" as never };
        const result = securityCheckResultSchema.safeParse(invalid);

        expect(result.success).toBe(false);
        if (!result.success) {
          expect(result.error.issues[0].message).toBe(
            "Please select a valid severity level"
          );
        }
      });

      it("should accept all valid severity levels", () => {
        const validSeverities = [
          "critical",
          "high",
          "medium",
          "low",
          "info",
        ] as const;

        validSeverities.forEach((severity) => {
          const valid = { ...validCheckResult, severity };
          const result = securityCheckResultSchema.safeParse(valid);
          expect(result.success).toBe(true);
        });
      });

      it("should reject invalid status values", () => {
        const invalid = { ...validCheckResult, status: "invalid" as never };
        const result = securityCheckResultSchema.safeParse(invalid);

        expect(result.success).toBe(false);
        if (!result.success) {
          expect(result.error.issues[0].message).toBe(
            "Please select a valid status"
          );
        }
      });

      it("should accept all valid status values", () => {
        const validStatuses = ["pass", "fail", "warning", "error"] as const;

        validStatuses.forEach((status) => {
          const valid = { ...validCheckResult, status };
          const result = securityCheckResultSchema.safeParse(valid);
          expect(result.success).toBe(true);
        });
      });
    });
  });

  describe("securityRuleSchema validation", () => {
    const validRule = {
      id: "rule-123",
      name: "Check Default Passwords",
      description: "Verify no default passwords are in use",
      vendor: "cisco" as const,
      command: "show running-config | include username",
      expectedPattern: "^(?!.*admin.*admin).*$",
      severity: "critical" as const,
      enabled: true,
      createdAt: new Date(),
    };

    it("should validate a valid security rule", () => {
      const result = securityRuleSchema.safeParse(validRule);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data).toEqual(validRule);
      }
    });

    it("should use default enabled value", () => {
      const ruleData = {
        id: validRule.id,
        name: validRule.name,
        description: validRule.description,
        vendor: validRule.vendor,
        command: validRule.command,
        expectedPattern: validRule.expectedPattern,
        severity: validRule.severity,
        createdAt: validRule.createdAt,
      };

      const result = securityRuleSchema.safeParse(ruleData);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data.enabled).toBe(true);
      }
    });

    it("should allow optional description and expectedPattern", () => {
      const minimalRule = {
        id: "rule-456",
        name: "Simple Check",
        vendor: "hp" as const,
        command: "show version",
        severity: "info" as const,
        createdAt: new Date(),
      };

      const result = securityRuleSchema.safeParse(minimalRule);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data.description).toBeUndefined();
        expect(result.data.expectedPattern).toBeUndefined();
        expect(result.data.enabled).toBe(true);
      }
    });

    describe("field validation", () => {
      it("should reject empty rule name", () => {
        const invalid = { ...validRule, name: "" };
        const result = securityRuleSchema.safeParse(invalid);

        expect(result.success).toBe(false);
        if (!result.success) {
          expect(result.error.issues[0].message).toBe("Rule name is required");
        }
      });

      it("should reject empty command", () => {
        const invalid = { ...validRule, command: "" };
        const result = securityRuleSchema.safeParse(invalid);

        expect(result.success).toBe(false);
        if (!result.success) {
          expect(result.error.issues[0].message).toBe("Command is required");
        }
      });

      it("should reject invalid vendor", () => {
        const invalid = { ...validRule, vendor: "invalid" as never };
        const result = securityRuleSchema.safeParse(invalid);

        expect(result.success).toBe(false);
        if (!result.success) {
          expect(result.error.issues[0].message).toBe(
            "Please select a valid vendor"
          );
        }
      });
    });
  });

  describe("securityCheckRequestSchema validation", () => {
    const validRequest = {
      deviceIds: ["device-1", "device-2"],
      ruleIds: ["rule-1", "rule-2"],
      parallel: true,
    };

    it("should validate a valid security check request", () => {
      const result = securityCheckRequestSchema.safeParse(validRequest);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data).toEqual(validRequest);
      }
    });

    it("should use default parallel value", () => {
      const requestData = {
        deviceIds: validRequest.deviceIds,
        ruleIds: validRequest.ruleIds,
      };

      const result = securityCheckRequestSchema.safeParse(requestData);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data.parallel).toBe(true);
      }
    });

    it("should allow optional ruleIds", () => {
      const requestWithoutRules = {
        deviceIds: ["device-1"],
      };

      const result = securityCheckRequestSchema.safeParse(requestWithoutRules);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data.ruleIds).toBeUndefined();
        expect(result.data.parallel).toBe(true);
      }
    });

    it("should reject empty deviceIds array", () => {
      const invalid = { ...validRequest, deviceIds: [] };
      const result = securityCheckRequestSchema.safeParse(invalid);

      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error.issues[0].message).toBe(
          "At least one device must be selected"
        );
      }
    });
  });

  describe("Security Check enums", () => {
    it("should have correct SecurityCheckStatus values", () => {
      expect(SecurityCheckStatus.PASS).toBe("pass");
      expect(SecurityCheckStatus.FAIL).toBe("fail");
      expect(SecurityCheckStatus.WARNING).toBe("warning");
      expect(SecurityCheckStatus.ERROR).toBe("error");
    });

    it("should have correct Severity values", () => {
      expect(Severity.CRITICAL).toBe("critical");
      expect(Severity.HIGH).toBe("high");
      expect(Severity.MEDIUM).toBe("medium");
      expect(Severity.LOW).toBe("low");
      expect(Severity.INFO).toBe("info");
    });
  });

  describe("SecurityCheckSummary interface", () => {
    it("should have correct structure", () => {
      const summary: SecurityCheckSummary = {
        deviceId: "device-123",
        deviceName: "Test Router",
        totalChecks: 10,
        passedChecks: 7,
        failedChecks: 2,
        warningChecks: 1,
        errorChecks: 0,
        lastChecked: new Date(),
        overallStatus: "warning",
      };

      expect(summary.deviceId).toBeDefined();
      expect(summary.deviceName).toBeDefined();
      expect(summary.totalChecks).toBe(10);
      expect(summary.passedChecks).toBe(7);
      expect(summary.failedChecks).toBe(2);
      expect(summary.warningChecks).toBe(1);
      expect(summary.errorChecks).toBe(0);
      expect(summary.overallStatus).toBe("warning");
    });
  });
});
