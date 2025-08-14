import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/devices")({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="p-6 animate-in slide-in-from-bottom-4 fade-in duration-300">
      <h2 className="text-2xl font-bold mb-4">Device Management</h2>
      <p className="text-muted-foreground">
        Manage your network devices - Coming Soon
      </p>
    </div>
  );
}
