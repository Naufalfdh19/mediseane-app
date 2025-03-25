import { AddressTypes, Response } from "../types/response";

export function ParseIDAndNameFromLocation(
  resp: Response<AddressTypes[]> | undefined
) {
  return resp?.data
    ? resp.data.map((item) => {
        return {
          id: item.id,
          name: item.name,
        };
      })
    : [];
}

export function GetLocationByID(
  id: string,
  resp: Response<AddressTypes[]> | undefined
) {
  return resp?.data.filter(function (item) {
    return item.id == id;
  })[0];
}
