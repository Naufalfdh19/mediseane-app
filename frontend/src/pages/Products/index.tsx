import { ArrowLeft, ArrowRight, MapPin } from "lucide-react";
import { useNavigate, useSearchParams } from "react-router-dom";

import ProductCard from "../../components/ProductCard";
import ErrorCard from "../../components/ErrorCard";
import SearchHeader from "../../components/SearchHeader";
import PageLoader from "../../components/PageLoader";

import useAxios from "../../hooks/useAxios";
import useTitle from "../../hooks/useTitle";
import {
  ProductResponse,
  Response,
} from "../../types/response";
import {
  API_METHOD_GET,
  API_USER_PRODUCTS,
  PATH_ADDRESS,
  PATH_LOGIN,
  PRODUCTS_TITLE,
  TOKEN_KEY,
} from "../../const/const";
import Button from "../../components/Button";

export default function Products() {
  const limitPerPage = 10;
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();

  useTitle(PRODUCTS_TITLE);

  const queryParams = new URLSearchParams(location.search);

  const search = queryParams.get("search");
  let page = queryParams.get("page");
  if (page === null) {
    page = "1"
  }
  const categoryDecode = queryParams.get("category")?.split("-");
  const categoryId = categoryDecode ? parseInt(categoryDecode[0]) : 0;
  const category = categoryDecode
    ? categoryDecode
        .filter(function (_, i) {
          return 0 !== i;
        })
        .join("-")
    : "";

  const {
    data: products,
    isLoading,
    error,
  } = useAxios<Response<ProductResponse>>(
    API_USER_PRODUCTS +
      "?" +
      (`page=${page}`) +
      `&limit=${limitPerPage}` +
      (search ? `&name=${search}` : "") +
      (categoryId ? `&category_id=${categoryId}` : ""),
    API_METHOD_GET
  );

  function changePage(page: number) {
    const newParams = new URLSearchParams(searchParams);
    newParams.set("page", page.toString());

    setSearchParams(newParams);
  }

  return (
    <>
      <SearchHeader />
      <div className="grow flex bg-primary-white justify-center">
        <div className="my-16 w-[90%] max-w-[1259px] flex flex-col items-center gap-10">
          <div className="w-full flex flex-col justify-between gap-4">
            <div className="flex flex-col gap-4">
              <h1 className="text-2xl font-bold text-primary-black">
                {search
                  ? `Search for: "${search}"`
                  : category
                  ? category
                  : "All Products"}
              </h1>
              <p className="text-primary-black font-bold">
                Showing results for:{" "}
                <span
                  className="cursor-pointer text-primary-gray font-normal"
                  onClick={() => {
                    if (localStorage.getItem(TOKEN_KEY)) {
                      navigate(PATH_ADDRESS);
                    } else {
                      navigate(PATH_LOGIN);
                    }
                  }}
                >
                  Your location <MapPin className="inline" />
                </span>
              </p>
            </div>
          </div>
          {error ? (
            <ErrorCard errors={error} />
          ) : isLoading ? (
            <PageLoader />
          ) : products?.data.products.length === 0 ? (
            <div className="w-full grow flex items-center justify-center">
              <div className="flex flex-col gap-2.5">
                <p className="text-primary-black font-semibold text-2xl">
                  No Results
                </p>
                <p className="text-primary-black">
                  Do you want to browse other products or{" "}
                  <span
                    className="cursor-pointer text-primary-blue font-bold"
                    onClick={() => {
                      if (localStorage.getItem(TOKEN_KEY)) {
                        navigate(PATH_ADDRESS);
                      } else {
                        navigate(PATH_LOGIN);
                      }
                    }}
                  >
                    change location?
                  </span>
                </p>
              </div>
            </div>
          ) : (
            <>
             <div className="w-full grid grid-cols-2 lg:grid-cols-5 justify-between gap-y-6 gap-x-2 md:justify-center md:gap-5">
          {products &&
            products.data &&
            products.data.products &&
            products.data.products.map((item) => {
              return <ProductCard key={item.id} product={item} />;
            })}
        </div>
        <div className="flex gap-2 mt-3">
          <div className="w-[50px]">
            <Button
              size="sm"
              square={true}
              onClick={() => {
                changePage(page ? +page - 1 : 1);
              }}
              disabled={page ? +page <= 1 : true}
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
                changePage(page ? +page + 1 : 1);
              }}
              disabled={products ? !(products.data.pagination.total_page > 1) : true}
            >
              <div className="flex justify-center items-center">
                <ArrowRight />
              </div>
            </Button>
          </div>
        </div>
            </>
          )}
        </div>
      </div>
    </>
  );
}
