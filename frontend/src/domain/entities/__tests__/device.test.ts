import { describe, it, expect } from "vitest";
import {
  deviceSchema,
  DeviceStatus,
  DeviceType,
  Vendor,
  type DeviceFormData,
  type Device,
} from "../device";

describe("Device Entity", () => {
  describe("deviceSchema validation", () => {
    const validDeviceData: DeviceFormData = {
      name: "Test Router",
      ipAddress: "192.168.1.1",
      deviceType: "router",
      vendor: "cisco",
      username: "admin",
      password: "password123",
      sshPort: 22,
      snmpCommunity: "public",
      tags: "production,core",
    };

    it("should validate a valid device", () => {
      const result = deviceSchema.safeParse(validDeviceData);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data).toEqual(validDeviceData);
      }
    });

    it("should use default SSH port when not provided", () => {
      const deviceData = {
        name: validDeviceData.name,
        ipAddress: validDeviceData.ipAddress,
        deviceType: validDeviceData.deviceType,
        vendor: validDeviceData.vendor,
        username: validDeviceData.username,
        password: validDeviceData.password,
        snmpCommunity: validDeviceData.snmpCommunity,
        tags: validDeviceData.tags,
      };

      const result = deviceSchema.safeParse(deviceData);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data.sshPort).toBe(22);
      }
    });

    it("should allow optional fields to be undefined", () => {
      const minimalDevice = {
        name: "Minimal Device",
        ipAddress: "10.0.0.1",
        deviceType: "switch" as const,
        vendor: "hp" as const,
        username: "user",
        password: "pass",
      };

      const result = deviceSchema.safeParse(minimalDevice);
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data.snmpCommunity).toBeUndefined();
        expect(result.data.tags).toBeUndefined();
        expect(result.data.sshPort).toBe(22);
      }
    });

    describe("field validation", () => {
      it("should reject empty device name", () => {
        const invalidDevice = { ...validDeviceData, name: "" };
        const result = deviceSchema.safeParse(invalidDevice);

        expect(result.success).toBe(false);
        if (!result.success) {
          expect(result.error.issues[0].message).toBe(
            "Device name is required"
          );
        }
      });

      it("should reject invalid IP addresses", () => {
        const invalidIPs = [
          "256.1.1.1",
          "192.168.1",
          "192.168.1.1.1",
          "not-an-ip",
          "192.168.1.256",
          "",
        ];

        invalidIPs.forEach((ip) => {
          const invalidDevice = { ...validDeviceData, ipAddress: ip };
          const result = deviceSchema.safeParse(invalidDevice);

          expect(result.success).toBe(false);
          if (!result.success) {
            expect(result.error.issues[0].message).toBe(
              "Invalid IP address format"
            );
          }
        });
      });

      it("should accept valid IP addresses", () => {
        const validIPs = [
          "0.0.0.0",
          "192.168.1.1",
          "10.0.0.1",
          "172.16.0.1",
          "255.255.255.255",
        ];

        validIPs.forEach((ip) => {
          const validDevice = { ...validDeviceData, ipAddress: ip };
          const result = deviceSchema.safeParse(validDevice);
          expect(result.success).toBe(true);
        });
      });

      it("should reject invalid device types", () => {
        const invalidDevice = {
          ...validDeviceData,
          deviceType: "invalid" as never,
        };
        const result = deviceSchema.safeParse(invalidDevice);

        expect(result.success).toBe(false);
        if (!result.success) {
          expect(result.error.issues[0].message).toBe(
            "Please select a valid device type"
          );
        }
      });

      it("should accept all valid device types", () => {
        const validTypes = [
          "router",
          "switch",
          "firewall",
          "access-point",
        ] as const;

        validTypes.forEach((type) => {
          const validDevice = { ...validDeviceData, deviceType: type };
          const result = deviceSchema.safeParse(validDevice);
          expect(result.success).toBe(true);
        });
      });

      it("should reject invalid vendors", () => {
        const invalidDevice = {
          ...validDeviceData,
          vendor: "invalid" as never,
        };
        const result = deviceSchema.safeParse(invalidDevice);

        expect(result.success).toBe(false);
        if (!result.success) {
          expect(result.error.issues[0].message).toBe(
            "Please select a valid vendor"
          );
        }
      });

      it("should accept all valid vendors", () => {
        const validVendors = ["cisco", "juniper", "hp", "aruba"] as const;

        validVendors.forEach((vendor) => {
          const validDevice = { ...validDeviceData, vendor };
          const result = deviceSchema.safeParse(validDevice);
          expect(result.success).toBe(true);
        });
      });

      it("should reject empty username", () => {
        const invalidDevice = { ...validDeviceData, username: "" };
        const result = deviceSchema.safeParse(invalidDevice);

        expect(result.success).toBe(false);
        if (!result.success) {
          expect(result.error.issues[0].message).toBe("Username is required");
        }
      });

      it("should reject empty password", () => {
        const invalidDevice = { ...validDeviceData, password: "" };
        const result = deviceSchema.safeParse(invalidDevice);

        expect(result.success).toBe(false);
        if (!result.success) {
          expect(result.error.issues[0].message).toBe("Password is required");
        }
      });

      it("should reject invalid SSH ports", () => {
        const invalidPorts = [0, -1, 65536, 100000];

        invalidPorts.forEach((port) => {
          const invalidDevice = { ...validDeviceData, sshPort: port };
          const result = deviceSchema.safeParse(invalidDevice);
          expect(result.success).toBe(false);
        });
      });

      it("should accept valid SSH ports", () => {
        const validPorts = [1, 22, 2222, 65535];

        validPorts.forEach((port) => {
          const validDevice = { ...validDeviceData, sshPort: port };
          const result = deviceSchema.safeParse(validDevice);
          expect(result.success).toBe(true);
        });
      });
    });
  });

  describe("Device enums", () => {
    it("should have correct DeviceStatus values", () => {
      expect(DeviceStatus.ONLINE).toBe("online");
      expect(DeviceStatus.OFFLINE).toBe("offline");
      expect(DeviceStatus.UNKNOWN).toBe("unknown");
    });

    it("should have correct DeviceType values", () => {
      expect(DeviceType.ROUTER).toBe("router");
      expect(DeviceType.SWITCH).toBe("switch");
      expect(DeviceType.FIREWALL).toBe("firewall");
      expect(DeviceType.ACCESS_POINT).toBe("access-point");
    });

    it("should have correct Vendor values", () => {
      expect(Vendor.CISCO).toBe("cisco");
      expect(Vendor.JUNIPER).toBe("juniper");
      expect(Vendor.HP).toBe("hp");
      expect(Vendor.ARUBA).toBe("aruba");
    });
  });

  describe("Device interface", () => {
    it("should extend DeviceFormData with system fields", () => {
      const device: Device = {
        id: "device-123",
        name: "Test Device",
        ipAddress: "192.168.1.1",
        deviceType: "router",
        vendor: "cisco",
        username: "admin",
        password: "password",
        sshPort: 22,
        status: "online",
        createdAt: new Date(),
        updatedAt: new Date(),
      };

      expect(device.id).toBeDefined();
      expect(device.status).toBeDefined();
      expect(device.createdAt).toBeDefined();
      expect(device.updatedAt).toBeDefined();
    });
  });
});
