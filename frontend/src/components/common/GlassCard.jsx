import Card from "./Card";

export default function GlassCard({ children, className = "", strong = false, ...props }) {
  const strengthClass = strong ? "glass-card-strong" : "";
  return (
    <Card className={`glass-card ${strengthClass} ${className}`.trim()} {...props}>
      {children}
    </Card>
  );
}
