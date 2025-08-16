import { z } from "zod";

// Device validation schema
export const deviceSchema = z.object({
  name: z.string().min(1, "Device name is required"),
  ipAddress: z
    .string()
    .regex(
      /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/,
      "Invalid IP address format"
    ),
  deviceType: z.enum(["router", "switch", "firewall", "access-point"], {
    message: "Please select a valid device type",
  }),
  vendor: z.enum(["cisco", "juniper", "hp", "aruba"], {
    message: "Please select a valid vendor",
  }),
  username: z.string().min(1, "Username is required"),
  password: z.string().min(1, "Password is required"),
  sshPort: z.number().min(1).max(65535).default(22),
  snmpCommunity: z.string().optional(),
  tags: z.string().optional(),
});

// Device form data type (for creating/editing devices)
export type DeviceFormData = z.infer<typeof deviceSchema>;

// Device entity type (includes system-generated fields)
export interface Device extends DeviceFormData {
  id: string;
  status: "online" | "offline" | "unknown";
  lastChecked?: Date;
  createdAt: Date;
  updatedAt: Date;
}

// Device status enum
export const DeviceStatus = {
  ONLINE: "online",
  OFFLINE: "offline",
  UNKNOWN: "unknown",
} as const;

export type DeviceStatusType = (typeof DeviceStatus)[keyof typeof DeviceStatus];

// Device type enum
export const DeviceType = {
  ROUTER: "router",
  SWITCH: "switch",
  FIREWALL: "firewall",
  ACCESS_POINT: "access-point",
} as const;

export type DeviceTypeType = (typeof DeviceType)[keyof typeof DeviceType];

// Vendor enum
export const Vendor = {
  CISCO: "cisco",
  JUNIPER: "juniper",
  HP: "hp",
  ARUBA: "aruba",
} as const;

export type VendorType = (typeof Vendor)[keyof typeof Vendor];
