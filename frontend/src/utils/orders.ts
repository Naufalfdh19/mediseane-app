import {
  ORDER_STATUS_CANCELED,
  ORDER_STATUS_DELIVERED,
  ORDER_STATUS_PENDING,
  ORDER_STATUS_PROCESSING,
  ORDER_STATUS_SHIPPED,
  ORDER_STATUS_WAITING,
} from "../const/const";

export function statusColor(status: string) {
  switch (status.toLowerCase()) {
    case ORDER_STATUS_WAITING:
    case ORDER_STATUS_PENDING:
      return "text-primary-gray";
    case ORDER_STATUS_PROCESSING:
      return "text-blue";
    case ORDER_STATUS_SHIPPED:
      return "text-primary-purple";
    case ORDER_STATUS_DELIVERED:
      return "text-primary-green";
    case ORDER_STATUS_CANCELED:
      return "text-primary-red";
    default:
      return "text-primary-black";
  }
}
