import { ShoppingCart } from "lucide-react";
import { useCartState } from "../../store/cart";

function CartIcon() {
  const cartState = useCartState();

  return (
    <div className="relative">
      <ShoppingCart />
      {cartState.items.length !== 0 && (
        <div className="w-[16px] h-[16px] absolute top-[-4px] right-[-6px] rounded-full bg-primary-red flex justify-center items-center">
          <p className="text-primary-white text-xs font-semibold text-center">
            {cartState.items.length}
          </p>
        </div>
      )}
    </div>
  );
}

export default CartIcon;
