import { GroupedCartItem } from "../../../types/response";
import CartItemBox from "./CartItem";

interface GroupedCartCardProps {
  items: GroupedCartItem;
  refetch: () => void
}

export default function GroupedCartCard({ items, refetch }: GroupedCartCardProps) {
  return (
    <div className="w-full bg-white shadow-md p-4 border border-primary-gray/20 rounded-lg flex flex-col gap-2">
      <div className="text-primary-black font-semibold">{items.pharmacy_name}</div>
      {items.items.map((item) => {
        return <CartItemBox key={item.cart_item_id} item={item} refetch={refetch}/>
      })}
    </div>
  );
}
