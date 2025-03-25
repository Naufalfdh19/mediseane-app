import { useEffect, useState } from "react";
import { CategoryData, Response } from "../../../types/response";
import useAxios from "../../../hooks/useAxios";
import { API_CATEGORIES, PATH_PRODUCTS } from "../../../const/const";
import ErrorCard from "../../../components/ErrorCard";
import {
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  ChevronUp,
} from "lucide-react";
import { useNavigate } from "react-router-dom";

import LoaderCard from "../../../components/LoaderCard";

export default function Categories() {
  const [pageCategory, setPageCategory] = useState(1);
  const [totalPageCategory, setTotalPageCategory] = useState(1);
  const [isCategoryHidden, setIsCategoryHidden] = useState<boolean>(true);

  const limitPerPageCategory = 8;

  const {
    data: categories,
    isLoading: isLoadingCategory,
    error: errCategory,
  } = useAxios<Response<CategoryData>>(
    API_CATEGORIES + `?page=${pageCategory}` + `&limit=${limitPerPageCategory}`
  );

  const navigate = useNavigate();

  useEffect(() => {
    const currPage =
      categories && categories.data && categories.data.pagination
        ? categories.data.pagination.current_page
        : 1;
    const totalPage =
      categories && categories.data && categories.data.pagination
        ? categories.data.pagination.total_page
        : 1;

    setTotalPageCategory(totalPage);
    setPageCategory(currPage);
  }, [categories]);

  useEffect(() => {
    if (pageCategory < 1) {
      setPageCategory(1);
    }
    if (pageCategory > totalPageCategory) {
      setPageCategory(totalPageCategory);
    }
  }, [pageCategory, totalPageCategory]);

  function handleBackButtonClick() {
    setPageCategory((page) => (page > 1 ? page - 1 : page));
  }

  function handleNextButtonClick() {
    setPageCategory((page) => (page < totalPageCategory ? page + 1 : page));
  }

  function renderCategories() {
    if (isLoadingCategory) {
      return (
        <div className="w-full h-full flex justify-center">
          <LoaderCard />
        </div>
      );
    }

    if (errCategory) {
      return <ErrorCard errors={errCategory} />;
    }

    if (
      !categories ||
      !categories.data ||
      !categories.data.list_item ||
      !(categories.data.list_item.length > 0)
    ) {
      return (
        <div className="w-full grow flex justify-center items-center">
          <p className="text-primary-gray">No categories found</p>
        </div>
      );
    }

    return (
      <>
        <div className="w-full flex justify-center items-center gap-[18px] mt-5 md:m-0 md:justify-between">
          <div className="w-[50px] hidden md:inline">
            <div
              className={`flex justify-center items-center ${
                pageCategory <= 1
                  ? "hidden"
                  : "text-primary-black cursor-pointer hover:text-primary-green"
              }`}
              onClick={handleBackButtonClick}
            >
              <ChevronLeft className="w-10 h-10" strokeWidth={0.8} />
            </div>
          </div>
          <div
            className={`w-full grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 md:justify-center gap-2.5 ${
              isCategoryHidden &&
              "h-[90px] md:h-[140px] xl:h-full overflow-hidden relative"
            }`}
          >
            {isCategoryHidden && (
              <div className="h-[50px] w-full absolute bottom-0 xl:bg-none xl:-z-10 bg-gradient-to-t from-primary-white to-primary-white/40 flex justify-center items-center rounded-lg"></div>
            )}
            {categories?.data.list_item.map((item) => {
              return (
                <div
                  key={item.id}
                  onClick={() => {
                    navigate(
                      PATH_PRODUCTS +
                        "?category=" +
                        item.id +
                        "-" +
                        encodeURIComponent(item.name)
                    );
                  }}
                  className="cursor-pointer"
                >
                  <div className="w-full bg-secondary-green px-5 py-2 rounded-full text-center">
                    <p className="text-primary-green font-semibold text-lg">{item.name}</p>
                  </div>
                </div>
              );
            })}
          </div>
          <div className="w-[50px] hidden md:inline">
            <div
              className={`flex justify-center items-center ${
                pageCategory >= totalPageCategory
                  ? "hidden"
                  : "text-primary-black cursor-pointer hover:text-primary-green"
              }`}
              onClick={handleNextButtonClick}
            >
              <ChevronRight className="w-10 h-10" strokeWidth={0.8} />
            </div>
          </div>
        </div>
        <div className="flex md:hidden mt-3">
          <div className="w-[50px]">
            <div
              className={`flex justify-center items-center ${
                pageCategory <= 1
                  ? "hidden"
                  : "text-primary-black cursor-pointer hover:text-primary-green"
              }`}
              onClick={handleBackButtonClick}
            >
              <ChevronLeft className="w-10 h-10" strokeWidth={0.8} />
            </div>
          </div>
          <div className="w-[50px]">
            <div
              className={`flex justify-center items-center ${
                pageCategory >= totalPageCategory
                  ? "hidden"
                  : "text-primary-black cursor-pointer hover:text-primary-green"
              }`}
              onClick={handleNextButtonClick}
            >
              <ChevronRight className="w-10 h-10" strokeWidth={0.8} />
            </div>
          </div>
        </div>
      </>
    );
  }

  return (
    <div className="w-full min-h-[89px] flex flex-col gap-4">
      <div className="flex items-center gap-2 justify-between">
        <p className="w-fit text-primary-black text-[26px] font-semibold overflow-auto">
          Categories
        </p>
        <div onClick={() => setIsCategoryHidden((prev) => !prev)}>
          {isCategoryHidden ? (
            <ChevronDown className="w-8 h-8 xl:hidden" strokeWidth={1.3} />
          ) : (
            <ChevronUp className="w-8 h-8 xl:hidden" strokeWidth={1.3} />
          )}
        </div>
      </div>
      <div>{renderCategories()}</div>
    </div>
  );
}
