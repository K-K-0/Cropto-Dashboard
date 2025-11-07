import { Icon, type LucideIcon } from "lucide-react";

interface StatsCardProps {
    icon: LucideIcon,
    value: number | string
    label: string
    variant: 'message' | 'trades' | 'tickers' | 'rate'
}

const StatsCard: React.FC<StatsCardProps> = ({ icon: Icon, value, label, variant}) => {
    const formatNumber = (num: number | string) => {
        if (typeof num == 'string') return num
        return new Intl.NumberFormat('en-US').format(num)
    } 

     return (
    <div className="stat-card">
      <div className={`stat-icon ${variant}`}>
        <Icon size={24} />
      </div>
      <div className="stat-content">
        <div className="stat-value">{formatNumber(value)}</div>
        <div className="stat-label">{label}</div>
      </div>
    </div>
  );
}

export default StatsCard