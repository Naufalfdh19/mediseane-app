import { useEffect, useState } from "react";

import Pagination from "../../../components/Pagination";
import TableHeader, {
  searchByType,
  sortByType,
} from "../../../components/TableHeader";
import PageLoader from "../../../components/PageLoader";

import useAxios from "../../../hooks/useAxios";
import { AdminProductData, Response } from "../../../types/response";
import { ADMIN_PRODUCTS_TITLE, API_ADMIN_PRODUCT } from "../../../const/const";
import useTitle from "../../../hooks/useTitle";

const searchByList: searchByType[] = [
  {
    id: "1",
    name: "name",
  },
  {
    id: "2",
    name: "generic_name",
  },
  {
    id: "3",
    name: "manufacture",
  },
  {
    id: "4",
    name: "product_classification",
  },
  {
    id: "5",
    name: "product_form",
  },
  {
    id: "6",
    name: "is_active",
  },
];

const sortByList: sortByType[] = [
  {
    id: "1",
    name: "created_at",
  },
  {
    id: "2",
    name: "name",
  },
  {
    id: "3",
    name: "product_used",
  },
];

export default function AdminProducts() {
  const limitPerPage = 10;
  const [search, setSearch] = useState("");
  const [page, setPage] = useState<number>(1);
  const [searchBy, setSearchBy] = useState<searchByType>(searchByList[0]);
  const [filter, setFilter] = useState("");
  const [sortBy, setSortBy] = useState<sortByType>();
  const [totalPage, setTotalPage] = useState<number>(1);
  const [asc, setAsc] = useState(false);

  const { data, isLoading } = useAxios<Response<AdminProductData>>(
    API_ADMIN_PRODUCT +
      `?page=${page}` +
      `&order=${asc ? "asc" : "desc"}` +
      (sortBy ? `&sort_by=${sortBy.name}` : "") +
      (searchBy && filter ? `&${searchBy.name}=${filter}` : "")
  );

  useTitle(ADMIN_PRODUCTS_TITLE);

  useEffect(() => {
    const currPage =
      data &&
      data.data &&
      data.data.pagination &&
      data.data.pagination.current_page >= 1
        ? data.data.pagination.current_page
        : 1;
    const totalPage =
      data &&
      data.data &&
      data.data.pagination &&
      data.data.pagination.total_page >= 1
        ? data.data.pagination.total_page
        : 1;

    setTotalPage(totalPage);
    setPage(currPage);
  }, [data]);

  useEffect(() => {
    if (page < 1) {
      setPage(1);
    }
    if (page > totalPage) {
      setPage(totalPage);
    }
  }, [page, totalPage]);

  if (isLoading) {
    return <PageLoader />;
  }

  return (
    <div className="grow flex bg-primary-white justify-center">
      <div className="my-16 w-[90%] max-w-[1108px] flex flex-col items-center gap-8">
        <div className="w-full flex justify-between gap-4">
          <h1 className="text-center text-2xl font-bold text-primary-black">
            Products
          </h1>
        </div>
        <TableHeader
          search={search}
          searchBy={searchBy}
          setSearchBy={setSearchBy}
          searchByList={searchByList}
          setSearch={setSearch}
          setFilter={setFilter}
          sortBy={sortBy}
          setSortBy={setSortBy}
          sortByList={sortByList}
          asc={asc}
          setAsc={setAsc}
          readOnly={true}
        />
        <div className="w-full flex flex-col gap-8 bg-clip-content overflow-x-auto">
          <table className="min-w-[627px]">
            <thead>
              <tr className="border-b border-b-primary-gray h-[58px] text-primary-black">
                <th className="w-[60px] px-4">No</th>
                <th className="text-left px-4">Name</th>
                <th className="text-left px-4">Generic Name</th>
                <th className="text-left px-4">Manufacturer</th>
                <th className="text-left px-4">Classification</th>
                <th className="text-left px-4">Form</th>
                <th className="text-left px-4">Active</th>
              </tr>
            </thead>
            <tbody>
              {data?.data && data.data.products && data.data.products.length > 0
                ? data.data.products.map((item, idx) => {
                    return (
                      <tr key={item.id} className="h-[58px] text-primary-black">
                        <td className="text-center w-[60px] px-4">
                          {limitPerPage * (page - 1) + idx + 1}
                        </td>
                        <td className="text-left px-4">{item.name}</td>
                        <td className="text-left px-4">{item.generic_name}</td>
                        <td className="text-left px-4">{item.manufacture}</td>
                        <td className="text-left px-4">
                          {item.product_classification}
                        </td>
                        <td className="text-left px-4">{item.product_form}</td>
                        <td className="text-left px-4">
                          {item.is_active ? "Yes" : "No"}
                        </td>
                      </tr>
                    );
                  })
                : undefined}
            </tbody>
          </table>
          {!(
            data?.data &&
            data.data.products &&
            data.data.products.length > 0
          ) && (
            <div className="w-full flex justify-center">
              <p className="text-primary-gray">Item Not Found</p>
            </div>
          )}
        </div>
        <div className="w-full flex justify-center md:justify-end">
          <Pagination page={page} setPage={setPage} totalPage={totalPage} />
        </div>
      </div>
    </div>
  );
}
