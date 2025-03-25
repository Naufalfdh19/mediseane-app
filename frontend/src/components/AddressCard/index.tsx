import Tag from "../Tag";

import { UserAddress } from "../../types/response";

interface AddressCardProps {
  address: UserAddress;
  onSelect?: () => void;
}

export default function AddressCard({
  address,
  onSelect = () => {},
}: AddressCardProps) {
  return (
    <div
      className="w-full bg-white cursor-pointer shadow-md p-4 border border-primary-gray/20 rounded-lg flex flex-col gap-2"
      onClick={onSelect}
    >
      <div className="w-full flex flex-wrap gap-3">
        <h1 className="text-primary-black font-semibold">{address.name}</h1>
        <h1 className="text-primary-gray">{address.phone_number}</h1>
      </div>
      <h1 className="text-primary-gray w-full">{address.address}</h1>
      <h1 className="text-primary-gray w-full">{`${address.province}, ${address.city}, ${address.district}, ${address.sub_district}, ${address.postal_code}`}</h1>
      {address.is_active && <Tag>Active</Tag>}
    </div>
  );
}
