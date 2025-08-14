import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/reports")({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="p-6 animate-in slide-in-from-bottom-4 fade-in duration-300">
      <h2 className="text-2xl font-bold mb-4">Reports</h2>
      <p className="text-muted-foreground">
        Generate compliance and security reports - Coming Soon
      </p>
    </div>
  );
}
