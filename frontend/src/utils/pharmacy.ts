import { PharmacyData, Response } from "../types/response";

export function ParseIDAndNameFromPharmacy(
  res: Response<PharmacyData> | undefined
) {
  return res?.data && res.data.pharmacies
    ? res.data.pharmacies.map((item) => {
        return {
          id: item.id.toString(),
          name: item.name,
        };
      })
    : [];
}

export function GetPharmacyByID(
  id: string,
  res: Response<PharmacyData> | undefined
) {
  return res?.data.pharmacies
    ? res?.data.pharmacies.filter(function (item) {
        return item.id.toString() == id;
      })[0]
    : undefined;
}
