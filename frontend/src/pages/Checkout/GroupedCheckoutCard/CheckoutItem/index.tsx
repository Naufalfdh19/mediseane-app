import { CartItem } from "../../../../types/response";

interface CheckoutItemProps {
  item: CartItem;
}

export default function CheckoutItem({ item }: CheckoutItemProps) {
  return (
    <div className="h-[150px] border-t-[1px] border-primary-gray/50 py-3 flex items-center overflow-auto gap-3 md:gap-6">
      <div className="flex gap-3 items-center grow">
        <div className="w-[60px] h-[60px] md:w-[80px] md:h-[80px] flex-shrink-0`">
          <img src={item.image} alt="image" />
        </div>
        <div className="grow">
          <div className="w-full max-w-[5000px] flex gap-2 justify-between">
            <div>
              <p className="truncate text-primary-black">{item.name}</p>
              <p className="truncate text-sm text-primary-gray">
                {item.manufacturer}
              </p>
            </div>
            <div className="bg-primary-green bg-opacity-[18%] min-w-[35px] h-[35px] p-1 flex items-center justify-center rounded">
              <p className="text-primary-green">x{item.quantity}</p>
            </div>
          </div>
          <div className="mt-7">
            <p className="text-primary-black">Total:</p>
            <p className="truncate text-primary-green font-semibold">
              Rp{item.subtotal}
            </p>
          </div>

          <div className="w-fit flex gap-3 items-center"></div>
        </div>
      </div>
    </div>
  );
}
