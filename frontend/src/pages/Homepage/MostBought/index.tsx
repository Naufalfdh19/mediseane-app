import { useEffect, useState } from "react";
import { ArrowLeft, ArrowRight } from "lucide-react";

import ErrorCard from "../../../components/ErrorCard";
import ProductCard from "../../../components/ProductCard";
import LoaderCard from "../../../components/LoaderCard";
import Button from "../../../components/Button";

import useAxios from "../../../hooks/useAxios";
import { ProductResponse, Response } from "../../../types/response";
import { API_USER_HOMEPAGE } from "../../../const/const";

export default function MostBought() {
  const [pageBought, setPageBought] = useState(1);
  const [totalPageBought, setTotalPageBought] = useState(1);

  const limitPerPageBoughts = 4;
  const {
    data: boughts,
    isLoading: isLoadingBoughts,
    error: errBoughts,
  } = useAxios<Response<ProductResponse>>(
    API_USER_HOMEPAGE + "?" + `page=${pageBought}` + `&limit=${limitPerPageBoughts}`
  );

  useEffect(() => {
    const currPage =
      boughts && boughts.data && boughts.data.pagination
        ? boughts.data.pagination.current_page
        : 1;
    const totalPage =
      boughts && boughts.data && boughts.data.pagination
        ? boughts.data.pagination.total_page
        : 1;

    setTotalPageBought(totalPage);
    setPageBought(currPage);
  }, [boughts]);

  useEffect(() => {
    if (pageBought < 1) {
      setPageBought(1);
    }
    if (pageBought > totalPageBought) {
      setPageBought(totalPageBought);
    }
  }, [pageBought, totalPageBought]);

  function renderProducts() {
    if (isLoadingBoughts) {
      return (
        <div className="w-full grow flex justify-center items-center">
          <LoaderCard />
        </div>
      );
    }

    if (errBoughts) {
      return <ErrorCard errors={errBoughts} />;
    }

    if (
      !boughts ||
      !boughts.data ||
      !boughts.data.products ||
      !(boughts.data.products.length > 0)
    ) {
      return (
        <div className="w-full grow flex justify-center items-center">
          <p className="text-primary-gray">Be the first to buy!</p>
        </div>
      );
    }

    return (
      <>
        <div className="w-full grid grid-cols-2 lg:grid-cols-4 justify-between gap-y-6 gap-x-2 md:justify-center md:gap-5">
          {boughts &&
            boughts.data &&
            boughts.data.products &&
            boughts.data.products.map((item) => {
              return <ProductCard key={item.id} product={item} />;
            })}
        </div>
        <div className="flex gap-2 mt-3">
          <div className="w-[50px]">
            <Button
              size="sm"
              square={true}
              onClick={() => {
                setPageBought((page) => (page > 1 ? page - 1 : page));
              }}
              disabled={pageBought <= 1}
            >
              <div className="flex justify-center items-center">
                <ArrowLeft />
              </div>
            </Button>
          </div>
          <div className="w-[50px]">
            <Button
              size="sm"
              square={true}
              onClick={() => {
                setPageBought((page) =>
                  page < totalPageBought ? page + 1 : page
                );
              }}
              disabled={pageBought >= totalPageBought}
            >
              <div className="flex justify-center items-center">
                <ArrowRight />
              </div>
            </Button>
          </div>
        </div>
      </>
    );
  }

  return (
    <div className="w-full min-h-[491px] flex flex-col items-center gap-4">
      <p className="w-full text-primary-black text-[26px] font-semibold overflow-auto">
        Most Bought Products
      </p>
      {renderProducts()}
    </div>
  );
}
