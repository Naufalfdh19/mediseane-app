import { ProductResponse, Response } from "../types/response";

export function ParseIDAndNameFromProduct(
  res: Response<ProductResponse> | undefined
) {
  return res?.data && res.data.products
    ? res.data.products.map((item) => {
        return {
          id: item.id.toString(),
          name: item.name,
        };
      })
    : [];
}

export function GetProductByID(
  id: string,
  res: Response<ProductResponse> | undefined
) {
  return res?.data.products
    ? res?.data.products.filter(function (item) {
        return item.id.toString() == id;
      })[0]
    : undefined;
}
