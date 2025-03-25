import { PartnersData, Response } from "../types/response";

export function ParseIDAndNameFromPartner(
  resp: Response<PartnersData> | undefined
) {
  return resp?.data && resp.data.partner_list
    ? resp.data.partner_list.map((item) => {
        return {
          id: item.id.toString(),
          name: item.name,
        };
      })
    : [];
}

export function GetPartnerByID(
  id: string,
  resp: Response<PartnersData> | undefined
) {
  return resp?.data.partner_list
    ? resp?.data.partner_list.filter(function (item) {
        return item.id.toString() == id;
      })[0]
    : undefined;
}
